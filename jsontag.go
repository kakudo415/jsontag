package jsontag

import (
	"go/ast"
	"go/token"
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
	Analyzer.Flags.IntVar(&where, "where", -1, "USAGE (TODO)")
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.Field)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		if n.Pos() != token.Pos(where) {
			return
		}
		pass.Report(analysis.Diagnostic{
			Pos:     n.Pos(),
			Message: "Want to add JSON tag like this? (-fix needed)\n" + fieldToJSONTag(n) + "\n",
			SuggestedFixes: []analysis.SuggestedFix{
				{
					Message: "Add Tag",
					TextEdits: []analysis.TextEdit{
						{
							Pos:     n.Pos(),
							End:     n.End(),
							NewText: []byte(fieldToJSONTag(n)),
						},
					},
				},
			},
		})
	})

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
