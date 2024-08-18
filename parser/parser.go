package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type StructField struct {
	Name string
	Type string
	Tag  string
}

type Index struct {
	Name    string
	Fields  []string
	Unique  bool
	Primary bool
}

type ParsedStruct struct {
	Name    string
	Fields  []StructField
	Indexes []Index
}

func ParseGormStruct(filePath string) (ParsedStruct, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return ParsedStruct{}, err
	}

	var parsedStruct ParsedStruct

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.TypeSpec:
			if _, ok := x.Type.(*ast.StructType); ok {
				parsedStruct.Name = x.Name.Name
			}
		case *ast.StructType:
			for _, field := range x.Fields.List {
				structField := StructField{
					Name: field.Names[0].Name,
					Type: typeToString(field.Type),
				}
				if field.Tag != nil {
					structField.Tag = strings.Trim(field.Tag.Value, "`")
					parsedStruct.Indexes = append(parsedStruct.Indexes, parseIndexes(structField.Name, structField.Tag)...)
				}
				parsedStruct.Fields = append(parsedStruct.Fields, structField)
			}
		}
		return true
	})

	return parsedStruct, nil
}

func typeToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + typeToString(t.X)
	case *ast.SelectorExpr:
		return typeToString(t.X) + "." + t.Sel.Name
	default:
		return ""
	}
}

func parseIndexes(fieldName, tag string) []Index {
	var indexes []Index
	gormTag := extractGormTag(tag)
	if gormTag == "" {
		return indexes
	}

	for _, option := range strings.Split(gormTag, ";") {
		option = strings.TrimSpace(option)
		if option == "index" {
			indexes = append(indexes, Index{Name: fieldName + "Index", Fields: []string{fieldName}, Unique: false})
		} else if option == "uniqueIndex" {
			indexes = append(indexes, Index{Name: fieldName + "UniqueIndex", Fields: []string{fieldName}, Unique: true})
		} else if strings.HasPrefix(option, "index:") {
			indexName := strings.TrimPrefix(option, "index:")
			indexes = append(indexes, Index{Name: indexName, Fields: []string{fieldName}, Unique: false})
		} else if strings.HasPrefix(option, "uniqueIndex:") {
			indexName := strings.TrimPrefix(option, "uniqueIndex:")
			indexes = append(indexes, Index{Name: indexName, Fields: []string{fieldName}, Unique: true})
		}
	}

	return indexes
}

func extractGormTag(tag string) string {
	for _, t := range strings.Split(tag, " ") {
		if strings.HasPrefix(t, "gorm:") {
			return strings.Trim(strings.TrimPrefix(t, "gorm:"), "\"")
		}
	}
	return ""
}
