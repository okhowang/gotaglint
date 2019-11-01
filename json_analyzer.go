package gotaglint

import (
	"fmt"
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/analysis"
	"strings"
	"unicode"
)

func isStringable(t types.Type) bool {
	switch b := t.(type) {
	case *types.Basic:
		switch b.Kind() {
		case types.Bool,
			types.Int, types.Int8, types.Int16, types.Int32, types.Int64,
			types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64, types.Uintptr,
			types.Float32, types.Float64,
			types.String:
			return true
		default:
			return false
		}
	case *types.Pointer:
		return isStringable(b.Elem())
	case *types.Named:
		return isStringable(b.Underlying())
	default:
		return false
	}
}

var JsonNameAnalyzer = analysis.Analyzer{
	Name: "jsontaglint",
	Doc:  "report unused json tag",
	Run: runTagChecker("json", func(pass *analysis.Pass, f *ast.Field, tag string) {
		jsonTag := ParseJsonTag(tag)
		if jsonTag.Skip {
			return
		}
		//check tag name
		if !isValidTag(jsonTag.Name) {
			pass.Reportf(f.Tag.Pos(), "invalid name:%q", render(pass.Fset, f.Tag))
		}
		//check for string option
		if jsonTag.String {
			if !isStringable(pass.TypesInfo.TypeOf(f.Type)) {
				pass.Report(analysis.Diagnostic{
					Pos:     f.Tag.Pos(),
					End:     f.Tag.End(),
					Message: fmt.Sprintf("string must use on scalar field:%q", render(pass.Fset, f.Tag)),
				})
			}
		}
		//check same name
		if len(f.Names) != 0 && jsonTag.Name == f.Names[0].String() {
			pass.Report(analysis.Diagnostic{
				Pos:     f.Tag.Pos(),
				End:     f.Tag.End(),
				Message: fmt.Sprintf("same json tag name %q", render(pass.Fset, f.Tag)),
			})
		}
	}),
}

type JsonTag struct {
	Name      string
	String    bool
	OmitEmpty bool
	Skip      bool
}

func ParseJsonTag(tag string) JsonTag {
	if tag == "-" {
		return JsonTag{Skip: true}
	}
	name, opts := parseTag(tag)
	var result JsonTag
	result.Name = name
	result.String = opts.Contains("string")
	result.OmitEmpty = opts.Contains("omitempty")
	return result
}

//Copy from encoding/json/encode.go

func isValidTag(s string) bool {
	for _, c := range s {
		switch {
		case strings.ContainsRune("!#$%&()*+-./:<=>?@[]^_{|}~ ", c):
			// Backslash and quote chars are reserved, but
			// otherwise any punctuation chars are allowed
			// in a tag name.
		case !unicode.IsLetter(c) && !unicode.IsDigit(c):
			return false
		}
	}
	return true
}

//Copy from encoding/json/tags.go

// tagOptions is the string following a comma in a struct field's "json"
// tag, or the empty string. It does not include the leading comma.
type tagOptions string

// parseTag splits a struct field's json tag into its name and
// comma-separated options.
func parseTag(tag string) (string, tagOptions) {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx], tagOptions(tag[idx+1:])
	}
	return tag, tagOptions("")
}

// Contains reports whether a comma-separated list of options
// contains a particular substr flag. substr must be surrounded by a
// string boundary or commas.
func (o tagOptions) Contains(optionName string) bool {
	if len(o) == 0 {
		return false
	}
	s := string(o)
	for s != "" {
		var next string
		i := strings.Index(s, ",")
		if i >= 0 {
			s, next = s[:i], s[i+1:]
		}
		if s == optionName {
			return true
		}
		s = next
	}
	return false
}
