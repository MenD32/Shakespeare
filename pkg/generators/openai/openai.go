package openai

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/MenD32/Shakespeare/pkg/huggingface"
	"github.com/MenD32/Shakespeare/pkg/trace"
)

const (
	DefaultEndpoint = "/v1/chat/completions"

	DefaultRole        = "system"
	DefaultContentType = "application/json"

	ContentHeader       = "Content-Type"
	AuthorizationHeader = "Authorization"
	AuthorizationPrefix = "Bearer "
)

type OpenAIGenerator struct {
	Config
	data huggingface.DatasetData
}

func NewOpenAIGenerator(conf Config) (*OpenAIGenerator, error) {
	generator := &OpenAIGenerator{
		Config: conf,
	}
	err := generator.loadDatasetData()
	if err != nil {
		return nil, fmt.Errorf("failed to load dataset data: %v", err)
	}

	return generator, nil
}

func (g *OpenAIGenerator) Generate() (trace.TraceLog, error) {

	requestCount := int(g.Config.GetRequestsPerSecond() * g.Config.GetDuration().Seconds())
	interval := time.Second / time.Duration(g.Config.GetRequestsPerSecond())
	traces := make([]trace.TraceLogRequest, 0)

	var body []byte

	method := "POST"
	headers := g.getHeaders()
	endpoint := g.getEndpoint()

	for i := 0; i < requestCount; i++ {
		prompt, err := g.getPrompt()
		if err != nil {
			return nil, fmt.Errorf("failed to get prompt: %v", err)
		}
		body, err = g.buildRequestBody(prompt)
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

func (g *OpenAIGenerator) loadDatasetData() error {
	client := huggingface.NewDefaultClient()
	parquetFiles, err := client.GetParquetFiles(g.Config.Dataset)
	if err != nil {
		return err
	}

	subset, ok := parquetFiles[g.Config.Split]
	if !ok {
		return fmt.Errorf("split %s not found", g.Config.Split)
	}
	parquetfileurls, ok := subset[g.Config.Subset]
	if !ok {
		return fmt.Errorf("subset %s not found", g.Config.Subset)
	}

	data := make(huggingface.DatasetData, 0)
	for _, url := range parquetfileurls {
		urldata, err := huggingface.NewDatasetData(url)
		if err != nil {
			return err
		}
		data = append(data, urldata...)
	}

	g.data = data
	return nil
}

func (g *OpenAIGenerator) getPrompt() (string, error) {
	prompt, err := huggingface.GetRandomPrompt(g.data, g.Config.Column, func(s string) bool {
		return len(s) >= g.Config.MinInputLength && len(s) <= g.Config.MaxInputLength
	})

	if err != nil {
		return "", fmt.Errorf("failed to get random prompt: %v", err)
	}

	return prompt, nil
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
