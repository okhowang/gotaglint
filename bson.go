package gotaglint

import (
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/analysis"
	"strings"
)

type BsonTag struct {
	Name           string
	Skip           bool
	OmitEmpty      bool
	MinSize        bool
	Inline         bool
	InvalidOptions []string
}

func ParseBsonTag(tag string) BsonTag {
	if tag == "-" {
		return BsonTag{Skip: true}
	}
	fields := strings.Split(tag, ",")
	var bsonTag BsonTag
	bsonTag.Name = fields[0]
	if len(fields) > 1 {
		for _, field := range fields[1:] {
			switch field {
			case "omitempty":
				bsonTag.OmitEmpty = true
			case "minsize":
				bsonTag.MinSize = true
			case "inline":
				bsonTag.Inline = true
			default:
				bsonTag.InvalidOptions = append(bsonTag.InvalidOptions, field)
			}
		}
	}
	return bsonTag
}

func isBsonInlineable(t types.Type) bool {
	switch b := t.(type) {
	case *types.Map:
		key, ok := b.Key().(*types.Basic)
		if !ok {
			return false
		}
		return key.Kind() == types.String
	case *types.Pointer:
		_, ok := b.Elem().(*types.Struct)
		return ok
	case *types.Struct:
		return true
	case *types.Named:
		return isBsonInlineable(b.Underlying())
	default:
		return false
	}
}

func isBsonMinSizeable(t types.Type) bool {
	switch b := t.(type) {
	case *types.Pointer:
		return isBsonMinSizeable(b.Elem())
	case *types.Named:
		return isBsonMinSizeable(b.Underlying())
	case *types.Basic:
		switch b.Kind() {
		case types.Uint64, types.Uintptr, types.Int64:
			return true
		default:
			return false
		}
	default:
		return false
	}
}

var BsonAnalyzer = analysis.Analyzer{
	Name: "bsontaglint",
	Doc:  "report bson tag problem",
	Run: runTagChecker("bson", func(pass *analysis.Pass, field *ast.Field, tag string) {
		bsonTag := ParseBsonTag(tag)
		if bsonTag.Skip {
			return
		}
		if len(bsonTag.InvalidOptions) > 0 {
			pass.Reportf(field.Tag.Pos(), "invalid bson options:%q", render(pass.Fset, field.Tag))
		}
		if len(field.Names) > 0 && bsonTag.Name == strings.ToLower(field.Names[0].Name) {
			pass.Reportf(field.Tag.Pos(), "same bson tag name:%q", render(pass.Fset, field.Tag))
		}
		if bsonTag.Inline && !isBsonInlineable(pass.TypesInfo.TypeOf(field.Type)) {
			pass.Reportf(field.Tag.Pos(), "inline must be struct, pointer to struct, map[string]*:%q", render(pass.Fset, field.Tag))
		}
		if bsonTag.MinSize && !isBsonMinSizeable(pass.TypesInfo.TypeOf(field.Type)) {
			pass.Reportf(field.Tag.Pos(), "minsize must be int64, uint64, uintptr:%q", render(pass.Fset, field.Tag))
		}
		//TODO check duplicate key and more than one map
	}),
}
