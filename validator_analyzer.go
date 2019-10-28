package gotaglint

import (
	"fmt"
	"github.com/fatih/structtag"
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/analysis"
	"strings"
)

var BindingNameAnalyzer = analysis.Analyzer{
	Name: "bindingtaglint",
	Doc:  "report bad usage for validator tag",
	Run:  runBindingTagLint,
}

func runBindingTagLint(pass *analysis.Pass) (interface{}, error) {
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
				if f == nil {
					continue
				}
				if f.Tag == nil {
					continue
				}
				tags, err := structtag.Parse(strings.Trim(f.Tag.Value, "`"))
				if err != nil {
					return true
				}
				for _, tag := range tags.Tags() {
					if tag.Key != "binding" {
						continue
					}
					if tag.Name == "exists" || tag.HasOption("exists") {
						t := pass.TypesInfo.TypeOf(f.Type)
						if t == nil {
							return true
						}
						switch t.Underlying().(type) {
						case *types.Pointer:
						case *types.Interface:
						default:
							pass.Report(analysis.Diagnostic{
								Pos:     f.Type.Pos(),
								End:     f.Type.End(),
								Message: fmt.Sprintf("exists field must be pointer/interface"),
							})
						}
					}
				}
			}
			return true
		})
	}
	return nil, nil
}
