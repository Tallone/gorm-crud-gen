package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Tallone/gorm-crud-gen/generator"
	"github.com/Tallone/gorm-crud-gen/parser"
	"github.com/spf13/cobra"
)

var (
	outputDir   string
	packageName string
	handler     bool
)

var rootCmd = &cobra.Command{
	Use:   "gorm-crud-generator",
	Short: "Generate CRUD services and handlers from GORM structs",
	Long:  `A CLI tool to generate CRUD services and Gin handlers with Swaggo comments from GORM structs.`,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringVarP(&outputDir, "output", "o", ".", "Output directory for generated files")
	generateCmd.Flags().StringVarP(&packageName, "package", "p", "main", "Package name for generated files")
	generateCmd.Flags().BoolVar(&handler, "handler", true, "Generate handler")
}

var generateCmd = &cobra.Command{
	Use:   "generate [input_file]",
	Short: "Generate CRUD code from a GORM struct",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		inputFile := args[0]

		// Parse the input file
		parsedStruct, err := parser.ParseGormStruct(inputFile)
		if err != nil {
			log.Panic(err)
		}

		// Create the generator
		gen := generator.NewGenerator(parsedStruct, packageName, outputDir, handler)

		// Generate the files
		gen.Generate()

		fmt.Printf("Successfully generated CRUD code for %s in %s\n", parsedStruct.Name, outputDir)
	},
}
