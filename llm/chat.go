package llm

import (
	"Mxx/llm/models"
	"Mxx/llm/runners"
	"fmt"
)

type Provider int

const (
	Ollama Provider = iota
	OpenAI
	Unknown
)

func PrepareRunner(provider Provider, opts map[string]string) (models.ChatRunner, error) {
	switch provider {
	case Ollama:
		return runners.GetOllamaRunner(opts)
	case OpenAI:
		return runners.GetOpenAIRunner(opts)
	}
	return nil, fmt.Errorf("invalid provider %v", provider)
}
