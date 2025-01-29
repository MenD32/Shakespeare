package config

import (
	"encoding/json"
	"fmt"

	"github.com/MenD32/Shakespeare/pkg/generators"
	"github.com/MenD32/Shakespeare/pkg/generators/openai"
)

type GeneratorType string

const (
	OpenAI GeneratorType = "openai"
)

func GeneratorFactory(generatorType GeneratorType, options string) (generators.Generator, error) {
	var err error
	switch generatorType {
	case OpenAI:
		openaiConfig := &openai.Config{}
		err = json.Unmarshal([]byte(options), openaiConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal OpenAI config: %v", err)
		}
		return openai.NewOpenAIGenerator(*openaiConfig)
	default:
		return nil, fmt.Errorf("unknown generator type: %s", generatorType)
	}
}
