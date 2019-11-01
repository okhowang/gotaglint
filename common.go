package gotaglint

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
	"golang.org/x/tools/go/analysis"
	"reflect"
	"strconv"
)

// render returns the pretty-print of the given node
func render(fset *token.FileSet, x interface{}) string {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, x); err != nil {
		panic(err)
	}
	return buf.String()
}

func runTagChecker(key string, fun func(analyzer *analysis.Pass, field *ast.Field, tag string)) func(pass *analysis.Pass) (interface{}, error) {
	return func(pass *analysis.Pass) (interface{}, error) {
		for _, file := range pass.Files {
			ast.Inspect(file, func(node ast.Node) bool {
				st, ok := node.(*ast.StructType)
				if !ok {
					return true
				}
				if st.Fields == nil {
					return true
				}
				for _, f := range st.Fields.List {
					if f.Tag == nil {
						continue
					}
					tv, err := strconv.Unquote(f.Tag.Value)
					if err != nil {
						pass.Reportf(f.Tag.Pos(), "invalid tag:%q", render(pass.Fset, f.Tag))
						continue
					}
					tags := reflect.StructTag(tv)
					tag, ok := tags.Lookup(key)
					if !ok {
						continue
					}
					fun(pass, f, tag)
				}
				return true
			})
		}
		return nil, nil
	}
}
