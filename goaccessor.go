package goaccessor

import (
	"go/ast"
	"log"
	"reflect"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var structAnalyzer = &analysis.Analyzer{
	Name: "goaccessor",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
	ResultType: reflect.TypeOf(StructMap{}),
}

type StructMap map[string]*ast.StructType

const doc = "goaccessor analyzer parses struct values which can be used for generating the accessors to guard null pointer access."

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.GenDecl)(nil),
	}

	var stractMap = StructMap{}
	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.GenDecl:
			for _, spec := range n.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok {
					return
				}
				if !ts.Name.IsExported() {
					return
				}
				st, ok := ts.Type.(*ast.StructType)
				if !ok {
					return
				}
				stractMap[ts.Name.String()] = st
			}
		}

		return
	})

	for tName, _ := range stractMap {
		log.Println(tName)
	}
	return stractMap, nil
}

var Analyzer = &analysis.Analyzer{
	Name: "goaccessor",
	Doc:  doc,
	Run:  run2,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
		structAnalyzer,
	},
}

func run2(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	smaps := pass.ResultOf[structAnalyzer].(StructMap)

	nodeFilter := []ast.Node{
		(*ast.GenDecl)(nil),
	}

	var stractMap = StructMap{}
	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.GenDecl:
			for _, spec := range n.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok {
					return
				}
				if !ts.Name.IsExported() {
					return
				}
				st, ok := ts.Type.(*ast.StructType)
				if !ok {
					return
				}
				stractMap[ts.Name.String()] = st
			}
		}

		return
	})

	for tName, _ := range stractMap {
		log.Println(tName)
	}
	return stractMap, nil
}
