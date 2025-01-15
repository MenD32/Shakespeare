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

	DefaultRole        = "system"
	DefaultContentType = "application/json"

	ContentHeader       = "Content-Type"
	AuthorizationHeader = "Authorization"
	AuthorizationPrefix = "Bearer "
)

type OpenAIGenerator struct {
	RPS      float64
	Duration time.Duration

	// OpenAI API key
	APIKey              string
	Method              string
	Endpoint            string
	Model               string
	MaxCompletionTokens int
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
		prompt := getPrompt()
		body, err = g.buildRequestBody(prompt[1])
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
		ContentHeader: DefaultContentType,
	}
	if g.APIKey != "" {
		headers[AuthorizationHeader] = AuthorizationPrefix + g.APIKey
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

func (g OpenAIGenerator) buildRequestBody(prompt string) ([]byte, error) {

	body := NewOpenAIRequestBody(
		g.Model,
		DefaultRole,
		prompt,
		g.MaxCompletionTokens,
	)

	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal body to JSON: %v", err)
	}

	return jsonData, nil
}

type OpenAIRequestBody struct {
	Model               string          `json:"model"`
	Messages            []OpenAIMessage `json:"messages"`
	Stream              bool            `json:"stream"`
	StreamOptions       StreamOptions   `json:"stream_options"`
	MaxCompletionTokens int             `json:"max_completion_tokens"`
	MaxTokens           int             `json:"max_tokens"` // Deprecated: Use MaxCompletionTokens instead
}

func NewOpenAIRequestBody(model string, role string, prompt string, maxCompletionTokens int) OpenAIRequestBody {
	return OpenAIRequestBody{
		Model: model,
		Messages: []OpenAIMessage{
			{
				Role:    role,
				Content: prompt,
			},
		},
		Stream: true,
		StreamOptions: StreamOptions{
			IncludeUsage: true,
		},
		MaxCompletionTokens: maxCompletionTokens,
	}
}

type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type StreamOptions struct {
	IncludeUsage bool `json:"include_usage"`
}
