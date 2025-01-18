package openai

type Config struct {
	APIKey   string
	Endpoint string
	Model    string

	MaxCompletionTokens int
}
