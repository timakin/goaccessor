package goaccessor

import (
	"go/ast"
	"log"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name: "goaccessor",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

const doc = "goaccessor analyzer parses struct values which can be used for generating the accessors to guard null pointer access."

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.GenDecl)(nil),
	}

	inspect.WithStack(nodeFilter, func(n ast.Node, push bool, stack []ast.Node) bool {
		switch n := n.(type) {
		case *ast.GenDecl:
			log.Printf("%+v", n.Specs)
			_ = n
		case *ast.Ident:
			_ = n
		}
		return false
	})

	return nil, nil
}
