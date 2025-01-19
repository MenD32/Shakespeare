package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/MenD32/Shakespeare/pkg/config"
)

var (
	outputFilePath string

	generatorType    string
	generatorOptions string // json

)

var rootCmd = &cobra.Command{
	Use:   "shakespeare",
	Short: "Shakespeare CLI application",
	Long:  `A CLI application to generate traces for the Tempest loadtesting tool.`,
	Run: func(cmd *cobra.Command, args []string) {

		generator, err := config.GeneratorFactory(config.GeneratorType(generatorType), generatorOptions)
		if err != nil {
			fmt.Printf("Error creating generator: %v\n", err)
			return
		}

		traceLog, err := generator.Generate()
		if err != nil {
			fmt.Printf("Error generating trace: %v\n", err)
			return
		}

		file, err := os.Create(outputFilePath)
		if err != nil {
			fmt.Printf("Error creating output file: %v\n", err)
			return
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		if err := encoder.Encode(traceLog); err != nil {
			fmt.Printf("Error encoding trace log to JSON: %v\n", err)
			return
		}
	},
}

func init() {
	rootCmd.Flags().StringVarP(&generatorType, "requestType", "t", "openai", "Type of request to generate")
	rootCmd.MarkFlagRequired("requestType")
	rootCmd.Flags().StringVarP(&generatorOptions, "generatorOptions", "g", "", "Options for the generator")
	rootCmd.MarkFlagRequired("generatorOptions")
	rootCmd.Flags().StringVarP(&outputFilePath, "output", "o", "trace.json", "Output file path")
	rootCmd.MarkFlagRequired("output")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
