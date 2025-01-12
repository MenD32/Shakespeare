package openai

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

const (
	DefaultMethod   = "POST"
	DefaultEndpoint = "/v1/chat/completions"
)

type OpenAIGenerator struct {
	RPS      float64
	Duration time.Duration

	// OpenAI API key
	APIKey   string
	Method   string
	Endpoint string
	Model    string
}

func (g *OpenAIGenerator) Generate() (trace.TraceLog, error) {

	requestCount := int(g.Duration.Seconds() * g.RPS)
	interval := time.Second / time.Duration(g.RPS)
	traces := make([]trace.TraceLogRequest, 0)

	var body []byte
	var err error

	headers := g.getHeaders()
	method := g.getMethod()
	endpoint := g.getEndpoint()

	for i := 0; i < requestCount; i++ {
		body, err = g.buildRequestbody()
		if err != nil {
			return nil, fmt.Errorf("failed to build request body: %v", err)
		}

		traces = append(traces, trace.TraceLogRequest{
			Delay:   time.Duration(i) * interval,
			Method:  method,
			Path:    endpoint,
			Headers: headers,
			Body:    body,
		})
	}
	return traces, nil

}

func (g *OpenAIGenerator) getMethod() string {
	if g.Method == "" {
		return DefaultMethod
	}
	return g.Method
}

func (g *OpenAIGenerator) getEndpoint() string {
	if g.Endpoint == "" {
		return DefaultEndpoint
	}
	return g.Endpoint
}

func (g *OpenAIGenerator) getHeaders() map[string]string {
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	if g.APIKey != "" {
		headers["Authorization"] = "Bearer " + g.APIKey
	}
	return headers
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

func (g OpenAIGenerator) buildRequestbody() ([]byte, error) {

	prompt := getPrompt()
	if len(prompt) != 2 {
		return nil, fmt.Errorf("prompt is Malformed")
	}

	body := make(map[string]interface{})
	body["model"] = g.Model
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

	body["max_completion_tokens"] = 300

	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal body to JSON: %v", err)
	}

	return jsonData, nil
}
