package generators

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/MenD32/Shakespeare/pkg/trace"
)

type OpenAIGenerator struct {
	RPS      float64
	Duration time.Duration
}

func (g *OpenAIGenerator) Generate() (trace.TraceLog, error) {

	requestCount := int(g.Duration.Seconds() * g.RPS)
	interval := time.Second / time.Duration(g.RPS)
	traces := make([]trace.TraceLogRequest, 0)

	var body []byte
	var err error

	for i := 0; i < requestCount; i++ {
		body, err = buildRequestbody()
		if err != nil {
			return nil, fmt.Errorf("failed to build request body: %v", err)
		}

		traces = append(traces, trace.TraceLogRequest{
			Delay:  time.Duration(i) * interval,
			Method: "POST",
			Path:   "/v1/chat/completions",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: body,
		})
	}
	return traces, nil

}

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

func getPrompt() []string {
	prompts := readCsvFile("/Users/mend/Downloads/prompts.csv")
	if len(prompts) == 0 {
		return []string{}
	}
	randomIndex := rand.Intn(len(prompts))
	return prompts[randomIndex]
}

func buildRequestbody() ([]byte, error) {

	prompt := getPrompt()
	if len(prompt) != 2 {
		return nil, fmt.Errorf("prompt is Malformed")
	}

	body := make(map[string]interface{})
	body["model"] = "Qwen/Qwen2-7B-Instruct"
	body["messages"] = []map[string]string{
		{
			"role":    "user",
			"content": prompt[1],
		},
	}
	body["stream"] = true
	body["stream_options"] = map[string]interface{}{
		"include_usage": true,
	}

	body["max_completion_tokens"] = 100

	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal body to JSON: %v", err)
	}

	return jsonData, nil
}
