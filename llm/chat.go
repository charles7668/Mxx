package llm

import (
	"Mxx/llm/models"
	"Mxx/llm/runners"
	"fmt"
)

type Provider int

const (
	Ollama Provider = iota
)

func PrepareRunner(provider Provider, opts map[string]string) (models.ChatRunner, error) {
	switch provider {
	case Ollama:
		return runners.GetOllamaRunner(opts)
	}
	return nil, fmt.Errorf("invalid provider %v", provider)
}
