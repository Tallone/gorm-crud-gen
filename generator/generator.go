package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Tallone/gorm-crud-gen/parser"
	"github.com/Tallone/gorm-crud-gen/templates"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
	// Ensure output directories exist
	dirs := []string{
		filepath.Join(g.OutputDir, "service"),
		filepath.Join(g.OutputDir, "handler"),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

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
		"title": cases.Title(language.English).String,
		"lower": strings.ToLower,
		"snake": toSnakeCase,
	}).ParseFS(templates.ServiceFile, "service.go.tmpl")
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
		"ServiceName": g.ParsedStruct.Name + "Service",
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

	outputPath := filepath.Join(g.OutputDir, "service", strings.ToLower(g.ParsedStruct.Name)+"_service.go")
	absolutePath, err := filepath.Abs(outputPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Check if file exists
	if _, err := os.Stat(absolutePath); os.IsNotExist(err) {
		// File doesn't exist, create it
		return os.WriteFile(absolutePath, formattedCode, 0644)
	}
	fmt.Println("Skip:", absolutePath)

	return nil
}

func (g *Generator) generateHandler() error {
	tmpl, err := template.New("handler.go.tmpl").Funcs(template.FuncMap{
		"title": cases.Title(language.English).String,
		"lower": strings.ToLower,
		"snake": toSnakeCase,
		"kebab": toKebabCase,
	}).ParseFS(templates.HandlerFile, "handler.go.tmpl")
	if err != nil {
		return err
	}

	fieldTypes := make(map[string]string)
	for _, field := range g.ParsedStruct.Fields {
		fieldTypes[field.Name] = field.Type
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{
		"StructName":   g.ParsedStruct.Name,
		"VarName":      strings.ToLower(g.ParsedStruct.Name[:1]) + g.ParsedStruct.Name[1:],
		"PackageName":  g.PackageName,
		"Indexes":      g.ParsedStruct.Indexes,
		"FieldTypes":   fieldTypes,
		"StructFields": g.ParsedStruct.Fields,
	})
	if err != nil {
		return err
	}

	formattedCode, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	outputPath := filepath.Join(g.OutputDir, "handler", strings.ToLower(g.ParsedStruct.Name)+"_handler.go")
	absolutePath, err := filepath.Abs(outputPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Check if file exists
	if _, err := os.Stat(absolutePath); os.IsNotExist(err) {
		// File doesn't exist, create it
		return os.WriteFile(absolutePath, formattedCode, 0644)
	}
	fmt.Println("Skip:", absolutePath)
	return err
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
