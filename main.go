package main

import (
	"fmt"
	"os"

	"github.com/Tallone/gorm-crud-gen/generator"
	"github.com/Tallone/gorm-crud-gen/parser"
	"github.com/spf13/cobra"
)

var (
	outputDir   string
	packageName string
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
			fmt.Printf("Error parsing input file: %v\n", err)
			os.Exit(1)
		}

		// Create the generator
		gen := generator.NewGenerator(parsedStruct, packageName, outputDir)

		// Generate the files
		if err := gen.Generate(); err != nil {
			fmt.Printf("Error generating files: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully generated CRUD code for %s in %s\n", parsedStruct.Name, outputDir)
	},
}
