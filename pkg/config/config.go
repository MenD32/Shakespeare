package config

type Config struct {
	RequestsPerSecond float64
	Duration          int
	OutputFilePath    string
}

type RequestType string

const (
	OpenAI RequestType = "openai"
)


