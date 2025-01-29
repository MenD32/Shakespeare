package openai

import "time"

const (
	DefaultAPIEndpoint       = "v1/chat/completions"
	DefaultMaxTokens         = 100
	DefaultRequestsPerSecond = 1
	DefaultDuration          = time.Minute
)

type Config struct {
	RequestsPerSecond float64 `json:"rps"`
	Duration          float64 `json:"duration"`

	APIKey              string `json:"apikey"`
	Endpoint            string `json:"endpoint"`
	Model               string `json:"model"`
	MaxCompletionTokens int    `json:"maxtokens"`

	MinInputLength int `json:"min_input_length"`
	MaxInputLength int `json:"max_input_length"`

	Dataset string `json:"dataset"`
	Split   string `json:"split"`
	Subset  string `json:"subset"`
	Column  string `json:"column"`
}

func (g *Config) GetRequestsPerSecond() float64 {
	return g.RequestsPerSecond
}

func (g *Config) GetDuration() time.Duration {
	if g.Duration == 0 {
		return DefaultDuration
	}
	return time.Second * time.Duration(g.Duration)
}

func (g *Config) GetAPIKey() string {
	return g.APIKey
}

func (g *Config) GetEndpoint() string {
	if g.Endpoint == "" {
		return DefaultAPIEndpoint
	}
	return g.Endpoint
}

func (g *Config) GetModel() string {
	return g.Model
}

func (g *Config) GetMaxCompletionTokens() int {
	if g.MaxCompletionTokens == 0 {
		return DefaultMaxTokens
	}
	return g.MaxCompletionTokens
}
