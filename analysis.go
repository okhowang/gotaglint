package gotaglint

import (
	"bytes"
	"fmt"
	"github.com/fatih/structtag"
	"go/ast"
	"go/printer"
	"go/token"
	"golang.org/x/tools/go/analysis"
	"strings"
)

// render returns the pretty-print of the given node
func render(fset *token.FileSet, x interface{}) string {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, x); err != nil {
		panic(err)
	}
	return buf.String()
}

func runJsonTagLint(pass *analysis.Pass) (interface{}, error) {
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
				if len(f.Names) == 0 {
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
					if tag.Key != "json" {
						continue
					}
					if tag.Name == f.Names[0].String() {
						if len(tag.Options) == 0 {
							tags.Delete("json")
						} else {
							tag.Name = ""
						}
						var newTags = ""
						if tags.Len() != 0 {
							newTags = fmt.Sprintf("`%s`", tags)
						}
						pass.Report(analysis.Diagnostic{
							Pos:     f.Tag.Pos(),
							End:     f.Tag.End(),
							Message: fmt.Sprintf("same json tag name %q", render(pass.Fset, f.Tag)),
							SuggestedFixes: []analysis.SuggestedFix{{
								Message: "",
								TextEdits: []analysis.TextEdit{{
									Pos:     f.Tag.Pos(),
									End:     f.Tag.End(),
									NewText: []byte(newTags),
								}},
							}},
						})
					}
				}
			}
			return true
		})
	}
	return nil, nil
}

var JsonNameAnalyzer = analysis.Analyzer{
	Name: "jsontaglint",
	Doc:  "report unused json tag",
	Run:  runJsonTagLint,
}
