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
			return nil, err
		}
		return openai.NewOpenAIGenerator(*openaiConfig), nil
	default:
		return nil, fmt.Errorf("unknown generator type: %s", generatorType)
	}
}
