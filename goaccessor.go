package goaccessor

import (
	"fmt"
	"go/ast"
	"log"
	"reflect"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

var Analyzer = &analysis.Analyzer{
	Name: "goaccessor",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
	ResultType: reflect.TypeOf(getters{}),
}

const (
	signYear = 2018

	fileSuffix = "-accessors.go"
)

var (
	// blacklistStruct lists structs to skip.
	blacklistStruct = map[string]bool{
		"Client":          true,
		"BasicCredential": true,
	}

	blacklistStructMethod = map[string]bool{
		"Client.GetBaseURL":    true,
		"Client.SetCredential": true,
	}
)

type StructMap map[string]*ast.StructType

const doc = "goaccessor analyzer parses struct values which can be used for generating the accessors to guard null pointer access."

func run(pass *analysis.Pass) (interface{}, error) {
	// inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// nodeFilter := []ast.Node{
	// 	// (*ast.)(nil),
	// 	(*ast.StructType)(nil),
	// }
	// inspect.Preorder(nodeFilter, func(n ast.Node) {
	// 	ts := pass.TypesInfo.Types[n.(*ast.StructType)].Type.(*types.Struct)
	// 	if !ts.Name.IsExported() {
	// 		return
	// 	}

	// 	// Check if the struct is blacklisted.
	// 	if blacklistStruct[ts.Name.Name] {
	// 		return
	// 	}

	// 	for _, field := range ts.Fields.List {
	// 		se, ok := field.Type.(*ast.StarExpr)
	// 		if len(field.Names) == 0 || !ok {
	// 			continue
	// 		}

	// 		fieldName := field.Names[0]
	// 		// Skip unexported identifiers.
	// 		if !fieldName.IsExported() {
	// 			continue
	// 		}
	// 		// Check if "struct.method" is blacklisted.
	// 		if key := fmt.Sprintf("%v.Get%v", ts.Name, fieldName); blacklistStructMethod[key] {
	// 			// logf("Method %v is blacklisted; skipping.", key)
	// 			continue
	// 		}

	// 		// log.Printf("%v", pass.TypesInfo.Defs[fieldName])
	// 		log.Printf("%+v", pass.TypesInfo.TypeOf(se))
	// 		log.Printf("%+v", se)
	// 		log.Printf("%+v", se.X)
	// 		log.Printf("%v", field.Names)
	// 		// log.Printf("%v", se.X)

	// 		// pass.TypesInfo.Defs[se.X.]
	// 		// switch x := se.X.(type) {
	// 		// case *ast.ArrayType:
	// 		// 	t.addArrayType(x, ts.Name.String(), fieldName.String())
	// 		// case *ast.Ident:
	// 		// 	t.addIdent(x, ts.Name.String(), fieldName.String())
	// 		// case *ast.MapType:
	// 		// 	t.addMapType(x, ts.Name.String(), fieldName.String())
	// 		// case *ast.SelectorExpr:
	// 		// 	t.addSelectorExpr(x, ts.Name.String(), fieldName.String())
	// 		// default:
	// 		// 	logf("processAST: type %q, field %q, unknown %T: %+v", ts.Name, fieldName, x, x)
	// 		// }
	// 	}

	// })

	for _, f := range pass.Files {
		for _, decl := range f.Decls {
			// log.Printf("%+v", pass.TypesInfo.TypeOf(decl))
			if decl, ok := decl.(*ast.GenDecl); ok {
				for _, spec := range decl.Specs {
					ts, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}
					if !ts.Name.IsExported() {
						continue
					}
					st, ok := ts.Type.(*ast.StructType)
					if !ok {
						continue
					}

					// Check if the struct is blacklisted.
					if blacklistStruct[ts.Name.Name] {
						continue
					}

					for _, field := range st.Fields.List {
						se, ok := field.Type.(*ast.StarExpr)
						if len(field.Names) == 0 || !ok {
							continue
						}

						fieldName := field.Names[0]
						// Skip unexported identifiers.
						if !fieldName.IsExported() {
							continue
						}
						// Check if "struct.method" is blacklisted.
						if key := fmt.Sprintf("%v.Get%v", ts.Name, fieldName); blacklistStructMethod[key] {
							// logf("Method %v is blacklisted; skipping.", key)
							continue
						}

						// log.Printf("%v", pass.TypesInfo.Defs[fieldName])
						log.Printf("%+v", pass.TypesInfo.TypeOf(se))
						log.Printf("%+v", se)
						log.Printf("%+v", se.X)
						log.Printf("%v", field.Names)
						// log.Printf("%v", se.X)

						// pass.TypesInfo.Defs[se.X.]
						// switch x := se.X.(type) {
						// case *ast.ArrayType:
						// 	t.addArrayType(x, ts.Name.String(), fieldName.String())
						// case *ast.Ident:
						// 	t.addIdent(x, ts.Name.String(), fieldName.String())
						// case *ast.MapType:
						// 	t.addMapType(x, ts.Name.String(), fieldName.String())
						// case *ast.SelectorExpr:
						// 	t.addSelectorExpr(x, ts.Name.String(), fieldName.String())
						// default:
						// 	logf("processAST: type %q, field %q, unknown %T: %+v", ts.Name, fieldName, x, x)
						// }
					}

					// pass.TypesInfo.Defs[st.]

					// log.Printf("%+v", st)
				}

				// 			//for _, spec := range decl.Specs {
				// 			//
				// 			//
				// 			//}

				// 			//if obj, ok := pass.TypesInfo.Defs[decl.Name].(*types.Func); ok {
				// 			//	pass.ExportObjectFact(obj, new(foundFact))
				// 			//}
			}
		}
	}

	return nil, nil
}

// Copyright 2018 The go-vrchat AUTHORS. All rights reserved.
////
//// Use of this source code is governed by a BSD-style
//// license that can be found in the LICENSE file.
//
//// +build ignore
//
//// gen-accessors generates accessor methods for structs with pointer fields.
////
//// It is meant to be used by the go-vrchat authors in conjunction with the
//// go generate tool before sending a commit to GitHub.
//package main
//
//import (
//"bytes"
//"flag"
//"fmt"
//"go/ast"
//"go/format"
//"go/parser"
//"go/token"
//"io/ioutil"
//"log"
//"os"
//"sort"
//"strings"
//"text/template"
//)
//
//var (
//	verbose = flag.Bool("v", false, "Print verbose log messages")
//
//	sourceTmpl = template.Must(template.New("source").Parse(source))
//
//	// blacklistStructMethod lists "struct.method" combos to skip.
//	blacklistStructMethod = map[string]bool{
//		"Client.GetBaseURL":         true,
//		"Client.SetCredential":      true,
//		"ErrorResponse.GetResponse": true,
//	}
//
//func logf(fmt string, args ...interface{}) {
//	if *verbose {
//		log.Printf(fmt, args...)
//	}
//}
//
//func main() {
//	flag.Parse()
//	fset := token.NewFileSet()
//
//	pkgs, err := parser.ParseDir(fset, ".", sourceFilter, 0)
//	if err != nil {
//		log.Fatal(err)
//		return
//	}
//
//	for pkgName, pkg := range pkgs {
//		t := &templateData{
//			filename: pkgName + fileSuffix,
//			Year:     signYear,
//			Package:  pkgName,
//			Imports:  map[string]string{},
//		}
//		for filename, f := range pkg.Files {
//			logf("Processing %v...", filename)
//			if err := t.processAST(f); err != nil {
//				log.Fatal(err)
//			}
//		}
//		if err := t.dump(); err != nil {
//			log.Fatal(err)
//		}
//	}
//	logf("Done.")
//}
//
//func (t *templateData) processAST(f *ast.File) error {
//	for _, decl := range f.Decls {
//		gd, ok := decl.(*ast.GenDecl)
//		if !ok {
//			continue
//		}
//		for _, spec := range gd.Specs {
//			ts, ok := spec.(*ast.TypeSpec)
//			if !ok {
//				continue
//			}
//			// Skip unexported identifiers.
//			if !ts.Name.IsExported() {
//				logf("Struct %v is unexported; skipping.", ts.Name)
//				continue
//			}
//			// Check if the struct is blacklisted.
//			if blacklistStruct[ts.Name.Name] {
//				logf("Struct %v is blacklisted; skipping.", ts.Name)
//				continue
//			}
//			st, ok := ts.Type.(*ast.StructType)
//			if !ok {
//				continue
//			}
//			for _, field := range st.Fields.List {
//				se, ok := field.Type.(*ast.StarExpr)
//				if len(field.Names) == 0 || !ok {
//					continue
//				}
//
//				fieldName := field.Names[0]
//				// Skip unexported identifiers.
//				if !fieldName.IsExported() {
//					logf("Field %v is unexported; skipping.", fieldName)
//					continue
//				}
//				// Check if "struct.method" is blacklisted.
//				if key := fmt.Sprintf("%v.Get%v", ts.Name, fieldName); blacklistStructMethod[key] {
//					logf("Method %v is blacklisted; skipping.", key)
//					continue
//				}
//
//				switch x := se.X.(type) {
//				case *ast.ArrayType:
//					t.addArrayType(x, ts.Name.String(), fieldName.String())
//				case *ast.Ident:
//					t.addIdent(x, ts.Name.String(), fieldName.String())
//				case *ast.MapType:
//					t.addMapType(x, ts.Name.String(), fieldName.String())
//				case *ast.SelectorExpr:
//					t.addSelectorExpr(x, ts.Name.String(), fieldName.String())
//				default:
//					logf("processAST: type %q, field %q, unknown %T: %+v", ts.Name, fieldName, x, x)
//				}
//			}
//		}
//	}
//	return nil
//}
//
//func sourceFilter(fi os.FileInfo) bool {
//	return !strings.HasSuffix(fi.Name(), "_test.go") && !strings.HasSuffix(fi.Name(), fileSuffix)
//}
//
//func (t *templateData) dump() error {
//	if len(t.Getters) == 0 {
//		logf("No getters for %v; skipping.", t.filename)
//		return nil
//	}
//
//	// Sort getters by ReceiverType.FieldName.
//	sort.Sort(byName(t.Getters))
//
//	var buf bytes.Buffer
//	if err := sourceTmpl.Execute(&buf, t); err != nil {
//		return err
//	}
//	clean, err := format.Source(buf.Bytes())
//	if err != nil {
//		return err
//	}
//
//	logf("Writing %v...", t.filename)
//	return ioutil.WriteFile(t.filename, clean, 0644)
//}
//
//func newGetter(receiverType, fieldName, fieldType, zeroValue string, namedStruct bool) *getter {
//	return &getter{
//		sortVal:      strings.ToLower(receiverType) + "." + strings.ToLower(fieldName),
//		ReceiverVar:  strings.ToLower(receiverType[:1]),
//		ReceiverType: receiverType,
//		FieldName:    fieldName,
//		FieldType:    fieldType,
//		ZeroValue:    zeroValue,
//		NamedStruct:  namedStruct,
//	}
//}
//
//func (t *templateData) addArrayType(x *ast.ArrayType, receiverType, fieldName string) {
//	var eltType string
//	switch elt := x.Elt.(type) {
//	case *ast.Ident:
//		eltType = elt.String()
//	default:
//		logf("addArrayType: type %q, field %q: unknown elt type: %T %+v; skipping.", receiverType, fieldName, elt, elt)
//		return
//	}
//
//	t.Getters = append(t.Getters, newGetter(receiverType, fieldName, "[]"+eltType, "nil", false))
//}
//
//func (t *templateData) addIdent(x *ast.Ident, receiverType, fieldName string) {
//	var zeroValue string
//	var namedStruct = false
//	switch x.String() {
//	case "int", "int64":
//		zeroValue = "0"
//	case "string":
//		zeroValue = `""`
//	case "bool":
//		zeroValue = "false"
//	case "Timestamp":
//		zeroValue = "Timestamp{}"
//	default:
//		zeroValue = "nil"
//		namedStruct = true
//	}
//
//	t.Getters = append(t.Getters, newGetter(receiverType, fieldName, x.String(), zeroValue, namedStruct))
//}
//
//func (t *templateData) addMapType(x *ast.MapType, receiverType, fieldName string) {
//	var keyType string
//	switch key := x.Key.(type) {
//	case *ast.Ident:
//		keyType = key.String()
//	default:
//		logf("addMapType: type %q, field %q: unknown key type: %T %+v; skipping.", receiverType, fieldName, key, key)
//		return
//	}
//
//	var valueType string
//	switch value := x.Value.(type) {
//	case *ast.Ident:
//		valueType = value.String()
//	default:
//		logf("addMapType: type %q, field %q: unknown value type: %T %+v; skipping.", receiverType, fieldName, value, value)
//		return
//	}
//
//	fieldType := fmt.Sprintf("map[%v]%v", keyType, valueType)
//	zeroValue := fmt.Sprintf("map[%v]%v{}", keyType, valueType)
//	t.Getters = append(t.Getters, newGetter(receiverType, fieldName, fieldType, zeroValue, false))
//}
//
//func (t *templateData) addSelectorExpr(x *ast.SelectorExpr, receiverType, fieldName string) {
//	if strings.ToLower(fieldName[:1]) == fieldName[:1] { // Non-exported field.
//		return
//	}
//
//	var xX string
//	if xx, ok := x.X.(*ast.Ident); ok {
//		xX = xx.String()
//	}
//
//	switch xX {
//	case "time", "json":
//		if xX == "json" {
//			t.Imports["encoding/json"] = "encoding/json"
//		} else {
//			t.Imports[xX] = xX
//		}
//		fieldType := fmt.Sprintf("%v.%v", xX, x.Sel.Name)
//		zeroValue := fmt.Sprintf("%v.%v{}", xX, x.Sel.Name)
//		if xX == "time" && x.Sel.Name == "Duration" {
//			zeroValue = "0"
//		}
//		t.Getters = append(t.Getters, newGetter(receiverType, fieldName, fieldType, zeroValue, false))
//	default:
//		logf("addSelectorExpr: xX %q, type %q, field %q: unknown x=%+v; skipping.", xX, receiverType, fieldName, x)
//	}
//}
//
//type templateData struct {
//	filename string
//	Year     int
//	Package  string
//	Imports  map[string]string
//	Getters  []*getter
//}
//
//
//const source = `// Copyright {{.Year}} The go-vrchat AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated by gen-accessors; DO NOT EDIT.

// package {{.Package}}
// {{with .Imports}}
// import (
//   {{- range . -}}
//   "{{.}}"
//   {{end -}}
// )
// {{end}}
// {{range .Getters}}
// {{if .NamedStruct}}
// // Get{{.FieldName}} returns the {{.FieldName}} field.
// func ({{.ReceiverVar}} *{{.ReceiverType}}) Get{{.FieldName}}() *{{.FieldType}} {
//   if {{.ReceiverVar}} == nil {
//     return {{.ZeroValue}}
//   }
//   return {{.ReceiverVar}}.{{.FieldName}}
// }
// {{else}}
// // Get{{.FieldName}} returns the {{.FieldName}} field if it's non-nil, zero value otherwise.
// func ({{.ReceiverVar}} *{{.ReceiverType}}) Get{{.FieldName}}() {{.FieldType}} {
//   if {{.ReceiverVar}} == nil || {{.ReceiverVar}}.{{.FieldName}} == nil {
//     return {{.ZeroValue}}
//   }
//   return *{{.ReceiverVar}}.{{.FieldName}}
// }
// {{end}}
// {{end}}
