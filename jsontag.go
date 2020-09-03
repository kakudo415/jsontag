package jsontag

import (
	"go/ast"
	"go/format"
	"go/token"
	"os"
	"strconv"
	"strings"

	"github.com/k0kubun/pp"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "jsontag is ..."

type intset map[int]bool

var where intset

func (i *intset) String() string {
	return pp.Sprintln(i)
}

func (i *intset) Set(v string) error {
	n, e := strconv.Atoi(v)
	if e != nil {
		return e
	}
	(*i)[n] = true
	return nil
}

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
	where = intset{}
	Analyzer.Flags.Var(&where, "where", "USAGE - TODO")
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.Field)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		if !where[int(n.Pos())] {
			return
		}
		switch n := n.(type) {
		case *ast.Field:
			n.Tag = &ast.BasicLit{}
			n.Tag.Kind = token.STRING
			n.Tag.Value = fieldToJSONTag(n)
		}
	})

	format.Node(os.Stdout, pass.Fset, pass.Files[0]) // FIXME ファイルが複数ないことを決め打ちしている

	return nil, nil
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
