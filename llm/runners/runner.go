package runners

import (
	"Mxx/llm/models"
	"fmt"
)

type Provider int

const (
	Ollama Provider = iota
	OpenAI
	Unknown
)

func PrepareRunner(provider Provider, optionsFunc ...func(options *models.RunnerOptions)) (models.ChatRunner, error) {
	opts := &models.RunnerOptions{}
	for _, optionFunc := range optionsFunc {
		if optionFunc == nil {
			continue
		}
		optionFunc(opts)
	}
	switch provider {
	case Ollama:
		return GetOllamaRunner(*opts)
	case OpenAI:
		return GetOpenAIRunner(*opts)
	default:
		return nil, fmt.Errorf("invalid provider %v", provider)
	}
}
