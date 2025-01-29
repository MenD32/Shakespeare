package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/MenD32/Shakespeare/pkg/config"
	"github.com/spf13/cobra"
	"k8s.io/klog"
)

var (
	outputFilePath   string
	generatorType    string
	generatorOptions string
)

var (
	Version        = "dev"
	CommitHash     = "none"
	BuildTimestamp = "unknown"
)

func BuildVersion() string {
	return fmt.Sprintf("%s-%s (%s)", Version, CommitHash, BuildTimestamp)
}

var runCmd = &cobra.Command{
	Use:   "generate",
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

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: `Get the version of Shakespeare.`,
	Long:  `Get the version of Shakespeare.`,
	Run: func(cmd *cobra.Command, args []string) {
		klog.Infof("Tempest version: %s\n", BuildVersion())
	},
}

var rootCmd = &cobra.Command{
	Short: "Shakespeare is a tool for creating traces for the Tempest loadtesting tool.",
	Long:  `Shakespeare is a tool for creating traces for the Tempest loadtesting tool.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	runCmd.Flags().StringVarP(&generatorType, "requestType", "t", "openai", "Type of request to generate")
	runCmd.MarkFlagRequired("requestType")
	runCmd.Flags().StringVarP(&generatorOptions, "generatorOptions", "g", "", "Options for the generator")
	runCmd.MarkFlagRequired("generatorOptions")
	runCmd.Flags().StringVarP(&outputFilePath, "output", "o", "trace.json", "Output file path")
	runCmd.MarkFlagRequired("output")

	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(versionCmd)
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
