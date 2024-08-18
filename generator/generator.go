package generator

import (
	"bytes"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Tallone/gorm-crud-gen/parser"
	"golang.org/x/text/cases"
)

type Generator struct {
	ParsedStruct parser.ParsedStruct
	PackageName  string
	OutputDir    string
}

func NewGenerator(parsedStruct parser.ParsedStruct, packageName, outputDir string) *Generator {
	return &Generator{
		ParsedStruct: parsedStruct,
		PackageName:  packageName,
		OutputDir:    outputDir,
	}
}

func (g *Generator) Generate() error {
	if err := g.generateService(); err != nil {
		return err
	}
	if err := g.generateHandler(); err != nil {
		return err
	}
	return nil
}

func (g *Generator) generateService() error {
	tmpl, err := template.New("service.go.tmpl").Funcs(template.FuncMap{
		"title": cases.Title,
		"lower": strings.ToLower,
		"snake": toSnakeCase,
	}).ParseFiles("templates/service.go.tmpl")
	if err != nil {
		return err
	}

	fieldTypes := make(map[string]string)
	for _, field := range g.ParsedStruct.Fields {
		fieldTypes[field.Name] = field.Type
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{
		"StructName":  g.ParsedStruct.Name,
		"VarName":     strings.ToLower(g.ParsedStruct.Name[:1]) + g.ParsedStruct.Name[1:],
		"PackageName": g.PackageName,
		"Indexes":     g.ParsedStruct.Indexes,
		"FieldTypes":  fieldTypes,
	})
	if err != nil {
		return err
	}

	formattedCode, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	outputPath := filepath.Join(g.OutputDir, "services", strings.ToLower(g.ParsedStruct.Name)+"_service.go")
	return os.WriteFile(outputPath, formattedCode, 0644)
}

func (g *Generator) generateHandler() error {
	tmpl, err := template.New("handler.go.tmpl").Funcs(template.FuncMap{
		"title": cases.Title,
		"lower": strings.ToLower,
		"snake": toSnakeCase,
		"kebab": toKebabCase,
	}).ParseFiles("templates/handler.go.tmpl")
	if err != nil {
		return err
	}

	fieldTypes := make(map[string]string)
	for _, field := range g.ParsedStruct.Fields {
		fieldTypes[field.Name] = field.Type
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{
		"StructName":  g.ParsedStruct.Name,
		"VarName":     strings.ToLower(g.ParsedStruct.Name[:1]) + g.ParsedStruct.Name[1:],
		"PackageName": g.PackageName,
		"Indexes":     g.ParsedStruct.Indexes,
		"FieldTypes":  fieldTypes,
	})
	if err != nil {
		return err
	}

	formattedCode, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	outputPath := filepath.Join(g.OutputDir, "handlers", strings.ToLower(g.ParsedStruct.Name)+"_handler.go")
	return os.WriteFile(outputPath, formattedCode, 0644)
}

func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

func toKebabCase(s string) string {
	return strings.ReplaceAll(toSnakeCase(s), "_", "-")
}
