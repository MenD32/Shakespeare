package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/MenD32/Shakespeare/pkg/generators"
)

var (
	requestsPerSecond int
	duration          int
	outputFilePath    string
)

var rootCmd = &cobra.Command{
	Use:   "shakespeare",
	Short: "Shakespeare CLI application",
	Long:  `A CLI application to generate traces for the Tempest loadtesting tool.`,
	Run: func(cmd *cobra.Command, args []string) {

		generator := generators.OpenAIGenerator{
			RPS:      float64(requestsPerSecond),
			Duration: time.Second * time.Duration(duration),
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
	rootCmd.Flags().IntVarP(&requestsPerSecond, "requestsPerSecond", "r", 0, "Number of requests per second")
	rootCmd.Flags().IntVarP(&duration, "duration", "d", 0, "Duration in seconds")
	rootCmd.Flags().StringVarP(&outputFilePath, "outputFilePath", "o", "", "Path to the output file")
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
