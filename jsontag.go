package jsontag

import (
	"errors"
	"go/ast"
	"go/format"
	"go/token"
	"os"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "jsontag is ..."

var where int

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
	Analyzer.Flags.IntVar(&where, "where", -1, "USAGE - TODO")
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

			n.Tag = &ast.BasicLit{}
			n.Tag.Kind = token.STRING
			n.Tag.Value = fieldToJSONTag(n)
			f := whoseChild(pass.Files, n)
			if f == nil {
				panic(errors.New("どのファイルにも属さないフィールド"))
			}
			format.Node(os.Stdout, pass.Fset, f)
		}
	})

	return nil, nil
}

func whoseChild(files []*ast.File, n ast.Node) *ast.File {
	for _, f := range files {
		if f.Pos() < n.Pos() && n.End() < f.End() {
			return f
		}
	}
	return nil
}

func fieldToJSONTag(n ast.Node) string {
	return "`json:\"" + toPascal(n.(*ast.Field).Names[0].Name) + "\"`"
}

func toPascal(a string) string {
	if len(a) < 2 {
		return a
	}
	return strings.ToLower(string([]rune(a)[0])) + string([]rune(a)[1:])
}
