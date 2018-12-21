package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"golang.org/x/tools/go/ast/inspector"
)

const (
	fileSuffix = "-accessors.go"
)

var (
	sourceTmpl = template.Must(template.New("source").Funcs(map[string]interface{}{
		"raw": func(text string) template.HTML {
			return template.HTML(text)
		},
	}).Parse(source))
)

func main() {
	flag.Parse()
	fset := token.NewFileSet()
	dirPath := "."
	if os.Args[1] != "" {
		dirPath = os.Args[1]
	}

	pkgs, err := parser.ParseDir(fset, dirPath, sourceFilter, 0)
	if err != nil {
		log.Fatal(err)
		return
	}

	for pkgName, pkg := range pkgs {
		d := &dumper{
			filename: pkgName + fileSuffix,
			Year:     time.Now().Year(),
			Package:  pkgName,
		}
		var files []*ast.File
		for _, f := range pkg.Files {
			files = append(files, f)
		}
		ai, err := ParseAccessors(files)
		if err != nil {
			log.Fatal(err)
		}
		d.Imports = ai.Imports
		d.Getters = ai.Getters
		if err := d.dump(); err != nil {
			log.Fatal(err)
		}
	}
}

func sourceFilter(fi os.FileInfo) bool {
	return !strings.HasSuffix(fi.Name(), "_test.go") && !strings.HasSuffix(fi.Name(), fileSuffix)
}

type getter struct {
	sortVal      string // Lower-case version of "ReceiverType.FieldName".
	ReceiverVar  string // The one-letter variable name to match the ReceiverType.
	ReceiverType string
	FieldName    string
	FieldType    string
	ZeroValue    string
	NamedStruct  bool // Getter for named struct.
}

type getters = []*getter

type accessorInfo struct {
	Imports map[string]string
	Getters getters
}

var (
	// blacklistStruct lists structs to skip.
	blacklistStruct = map[string]bool{
		"Client":          true,
		"BasicCredential": true,
	}

	blacklistStructMethod = map[string]bool{
		"User.GetUnsubscribe": true,
	}
)

// ParseAccessors parses struct values which can be used for generating the accessors to guard null pointer access."
func ParseAccessors(files []*ast.File) (*accessorInfo, error) {
	inspect := inspector.New(files)

	nodeFilter := []ast.Node{
		(*ast.TypeSpec)(nil),
	}

	ai := &accessorInfo{
		Getters: []*getter{},
		Imports: map[string]string{},
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch ts := n.(type) {
		case *ast.TypeSpec:
			if !ts.Name.IsExported() {
				return
			}

			st, ok := ts.Type.(*ast.StructType)
			if !ok {
				return
			}

			// Check if the struct is blacklisted.
			if blacklistStruct[ts.Name.Name] {
				return
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
				key := fmt.Sprintf("%v.Get%v", ts.Name, fieldName)
				log.Println(key)
				if blacklistStructMethod[key] {
					continue
				}

				switch x := se.X.(type) {
				case *ast.ArrayType:
					addArrayType(ai, x, ts.Name.String(), fieldName.String())
				case *ast.Ident:
					addIdent(ai, x, ts.Name.String(), fieldName.String())
				case *ast.MapType:
					addMapType(ai, x, ts.Name.String(), fieldName.String())
				case *ast.SelectorExpr:
					addSelectorExpr(ai, x, ts.Name.String(), fieldName.String())
				}
			}

		}
		return
	})

	return ai, nil
}

func newGetter(receiverType, fieldName, fieldType, zeroValue string, namedStruct bool) *getter {
	return &getter{
		sortVal:      strings.ToLower(receiverType) + "." + strings.ToLower(fieldName),
		ReceiverVar:  strings.ToLower(receiverType[:1]),
		ReceiverType: receiverType,
		FieldName:    fieldName,
		FieldType:    fieldType,
		ZeroValue:    zeroValue,
		NamedStruct:  namedStruct,
	}
}

func addArrayType(ai *accessorInfo, x *ast.ArrayType, receiverType, fieldName string) {
	var eltType string
	switch elt := x.Elt.(type) {
	case *ast.Ident:
		eltType = elt.String()
	default:
		return
	}

	ai.Getters = append(ai.Getters, newGetter(receiverType, fieldName, "[]"+eltType, "nil", false))
}

func addIdent(ai *accessorInfo, x *ast.Ident, receiverType, fieldName string) {
	var zeroValue string
	var namedStruct = false
	switch x.String() {
	case "int", "int64":
		zeroValue = "0"
	case "string":
		zeroValue = `""`
	case "bool":
		zeroValue = "false"
	case "Timestamp":
		zeroValue = "Timestamp{}"
	default:
		zeroValue = "nil"
		namedStruct = true
	}

	ai.Getters = append(ai.Getters, newGetter(receiverType, fieldName, x.String(), zeroValue, namedStruct))
}

func addMapType(ai *accessorInfo, x *ast.MapType, receiverType, fieldName string) {
	var keyType string
	switch key := x.Key.(type) {
	case *ast.Ident:
		keyType = key.String()
		return
	}

	var valueType string
	switch value := x.Value.(type) {
	case *ast.Ident:
		valueType = value.String()
		return
	}

	fieldType := fmt.Sprintf("map[%v]%v", keyType, valueType)
	zeroValue := fmt.Sprintf("map[%v]%v{}", keyType, valueType)
	ai.Getters = append(ai.Getters, newGetter(receiverType, fieldName, fieldType, zeroValue, false))
}

func addSelectorExpr(ai *accessorInfo, x *ast.SelectorExpr, receiverType, fieldName string) {
	if strings.ToLower(fieldName[:1]) == fieldName[:1] { // Non-exported field.
		return
	}

	var xX string
	if xx, ok := x.X.(*ast.Ident); ok {
		xX = xx.String()
	}

	switch xX {
	case "time", "json":
		if xX == "json" {
			ai.Imports["encoding/json"] = "encoding/json"
		} else {
			ai.Imports[xX] = xX
		}
		fieldType := fmt.Sprintf("%v.%v", xX, x.Sel.Name)
		zeroValue := fmt.Sprintf("%v.%v{}", xX, x.Sel.Name)
		if xX == "time" && x.Sel.Name == "Duration" {
			zeroValue = "0"
		}
		ai.Getters = append(ai.Getters, newGetter(receiverType, fieldName, fieldType, zeroValue, false))
	}
}

type dumper struct {
	accessorInfo
	filename string
	Year     int
	Package  string
}

func (t *dumper) dump() error {
	if len(t.Getters) == 0 {
		return nil
	}

	// Sort getters by ReceiverType.FieldName.
	sort.Slice(t.Getters, func(i, j int) bool {
		//log.Printf("%+v", t.Getters[i])
		return t.Getters[i].sortVal < t.Getters[j].sortVal
	})

	var buf bytes.Buffer
	if err := sourceTmpl.Execute(&buf, t); err != nil {
		return err
	}
	clean, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	return ioutil.WriteFile(t.filename, clean, 0644)
}

const source = `// Copyright {{.Year}} The AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated by github.com/timakin/goaccessor; DO NOT EDIT.

package {{.Package}}
{{with .Imports}}
import (
  {{- range . -}}
  "{{.}}"
  {{end -}}
)
{{end}}
{{range .Getters}}
{{if .NamedStruct}}
// Get{{.FieldName}} returns the {{.FieldName}} field.
func ({{.ReceiverVar}} *{{.ReceiverType}}) Get{{.FieldName}}() *{{.FieldType}} {
  if {{.ReceiverVar}} == nil {
    return {{.ZeroValue | raw}}
  }
  return {{.ReceiverVar}}.{{.FieldName}}
}
{{else}}
// Get{{.FieldName}} returns the {{.FieldName}} field if it's non-nil, zero value otherwise.
func ({{.ReceiverVar}} *{{.ReceiverType}}) Get{{.FieldName}}() {{.FieldType}} {
  if {{.ReceiverVar}} == nil || {{.ReceiverVar}}.{{.FieldName}} == nil {
    return {{.ZeroValue | raw}}
  }
  return *{{.ReceiverVar}}.{{.FieldName}}
}
{{end}}
{{end}}
`
