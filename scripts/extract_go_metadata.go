package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

type Metadata struct {
	Package string   `json:"package"`
	Exports []string `json:"exports"`
	Imports []string `json:"imports"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <file.go>\n", os.Args[0])
		os.Exit(1)
	}

	filename := os.Args[1]
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing file: %v\n", err)
		os.Exit(1)
	}

	meta := Metadata{
		Package: f.Name.Name,
		Exports: []string{},
		Imports: []string{},
	}

	for _, imp := range f.Imports {
		val := imp.Path.Value
		if len(val) >= 2 && val[0] == '"' && val[len(val)-1] == '"' {
			val = val[1 : len(val)-1]
		}
		meta.Imports = append(meta.Imports, val)
	}

	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			if d.Name.IsExported() {
				meta.Exports = append(meta.Exports, d.Name.Name)
			}
		case *ast.GenDecl:
			for _, spec := range d.Specs {
				switch s := spec.(type) {
				case *ast.TypeSpec:
					if s.Name.IsExported() {
						meta.Exports = append(meta.Exports, s.Name.Name)
					}
				case *ast.ValueSpec:
					for _, name := range s.Names {
						if name.IsExported() {
							meta.Exports = append(meta.Exports, name.Name)
						}
					}
				}
			}
		}
	}

	out, err := json.Marshal(meta)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error formatting JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(out))
}
