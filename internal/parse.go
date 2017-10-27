package internal

import (
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"

	"github.com/lazada/sqle"
)

func ParseFile(path, tag string, nam sqle.NamingConvention) (fi *FileInfo, err error) {
	var (
		file *ast.File
		decl *ast.GenDecl
		typ  *ast.TypeSpec
		st   *ast.StructType
		si   *StructInfo
		ok   bool
	)
	if file, err = parser.ParseFile(token.NewFileSet(), path, nil, parser.ParseComments); err != nil {
		return
	}
	fi = &FileInfo{Package: file.Name.Name}

	for _, d := range file.Decls {
		if decl, ok = d.(*ast.GenDecl); !ok {
			continue
		}
		for _, s := range decl.Specs {
			if typ, ok = s.(*ast.TypeSpec); !ok {
				continue
			}
			if st, ok = typ.Type.(*ast.StructType); !ok || st.Incomplete {
				continue
			}
			if typ.Doc == nil {
				typ.Doc = decl.Doc
			}
			si, ok = &StructInfo{Name: typ.Name.Name}, false

			for _, f := range st.Fields.List {
				if len(f.Names) == 0 || !f.Names[0].IsExported() {
					continue
				}
				if f.Tag != nil {
					if path = reflect.StructTag(f.Tag.Value[1:]).Get(tag); path == `-` {
						continue
					}
					if path == `` {
						path = nam.Name(f.Names[0].Name)
					}
					si.Field = append(si.Field, &FieldInfo{Name: f.Names[0].Name, Alias: path})
					ok = true
				} else if tag == `` {
					si.Field = append(si.Field, &FieldInfo{Name: f.Names[0].Name, Alias: nam.Name(f.Names[0].Name)})
					ok = true
				}
			}
			if ok {
				fi.Struct = append(fi.Struct, si)
			}
		}
	}
	return
}
