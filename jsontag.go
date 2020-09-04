package jsontag

import (
	"errors"
	"go/ast"
	"go/format"
	"go/token"
	"os"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "jsontag is ..."

var where int
var option string

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "jsontag",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func init() {
	Analyzer.Flags.IntVar(&where, "where", -1, "タグを追加したいフィールドのオフセット位置")
	Analyzer.Flags.StringVar(&option, "option", "", "タグのオプション、詳しくは README.md を読んでください")
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.Field)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.Field:
			if token.Pos(where) < n.Pos() || n.End() < token.Pos(where) {
				return
			}

			if !n.Names[0].IsExported() {
				pass.Reportf(n.Pos(), n.Names[0].Name+" is NOT exported")
				return
			}

			if n.Tag == nil {
				n.Tag = &ast.BasicLit{}
				n.Tag.Kind = token.STRING
			}

			tags := parseFieldTag(n.Tag.Value)
			var jsontag *fieldTag
			for _, tag := range tags {
				if tag.key == "json" {
					jsontag = tag
				}
			}

			if jsontag != nil {
				return
			}

			jsontag = &fieldTag{
				key:   "json",
				value: fieldToJSONTagValue(n, option),
			}

			tags = append(tags, jsontag)

			n.Tag.Value = formatFieldTag(tags)

			f := whoseChild(pass.Files, n)
			if f == nil {
				panic(errors.New("どのファイルにも属さないフィールド"))
			}
			file, err := os.OpenFile(pass.Fset.File(n.Pos()).Name(), os.O_WRONLY, 0666)
			if err != nil {
				panic(err)
			}
			format.Node(file, pass.Fset, f)
		}
	})

	return nil, nil
}

type fieldTag struct {
	key   string
	value string
}

func parseFieldTag(src string) []*fieldTag {
	tags := []*fieldTag{}
	src, err := strconv.Unquote(src)
	if err != nil {
		return nil
	}

	tokens := tagLexer(src)

	for i := 0; i < len(tokens); {
		tag := fieldTag{}
		if (i+1 < len(tokens)) && tokens[i+1] == ":" {
			tag.key = tokens[i]
			tag.value = tokens[i+2]
			i += 3
		} else {
			tag.key = tokens[i]
			i++
		}
		tags = append(tags, &tag)
	}

	return tags
}

func tagLexer(src string) []string {
	var tokens []string
	chars := []rune(src)

	for pos := 0; pos < len(chars); {
		switch {
		case unicode.IsLetter(chars[pos]):
			var token string
			for pos < len(chars) && unicode.IsLetter(chars[pos]) {
				token += string(chars[pos])
				pos++
			}
			tokens = append(tokens, token)
		case chars[pos] == ':':
			tokens = append(tokens, ":")
			pos++
		case chars[pos] == '"':
			pos++
			var token string
			for pos < len(chars) && unicode.IsLetter(chars[pos]) {
				token += string(chars[pos])
				pos++
			}
			tokens = append(tokens, token)
			pos++
		case unicode.IsSpace(chars[pos]):
			for pos < len(chars) && unicode.IsSpace(chars[pos]) {
				pos++
			}
		}
	}
	return tokens
}

func formatFieldTag(tags []*fieldTag) string {
	s := "`"
	for i, tag := range tags {
		if i > 0 {
			s += " "
		}
		if len(tag.value) > 0 {
			s += tag.key + `:"` + tag.value + `"`
		} else {
			s += tag.key
		}
	}
	s += "`"
	return s
}

func whoseChild(files []*ast.File, n ast.Node) *ast.File {
	for _, f := range files {
		if f.Pos() < n.Pos() && n.Pos() < f.End() {
			return f
		}
	}
	return nil
}

func fieldToJSONTagValue(n ast.Node, o string) string {
	if o == "omitempty" {
		return toCamel(n.(*ast.Field).Names[0].Name) + ",omitempty"
	}
	if o == "ignore" {
		return "-"
	}
	return toCamel(n.(*ast.Field).Names[0].Name)
}

func toCamel(a string) string {
	if len(a) < 2 {
		return a
	}
	return strings.ToLower(string([]rune(a)[0])) + string([]rune(a)[1:])
}
