package runners

import (
	"Mxx/llm/models"
	"context"
	"fmt"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"log"
	"strings"
)

type OllamaRunner struct {
	models.ChatRunner
	modelName string
}

func GetOllamaRunner(options map[string]string) (*OllamaRunner, error) {
	result := &OllamaRunner{}
	if _, exist := options["model"]; !exist {
		return nil, fmt.Errorf("model not found in options")
	}
	return result, nil
}

func (r *OllamaRunner) Chat(ctx context.Context, characterMessage, prompt, contextMessage string) (string, error) {
	llm, err := ollama.New(ollama.WithModel("mistral"))
	if err != nil {
		log.Fatal(err)
	}
	llmCtx, cancelLLM := context.WithCancel(ctx)

	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeHuman, prompt),
	}
	if characterMessage != "" {
		content = append(content, llms.TextParts(llms.ChatMessageTypeSystem, characterMessage))
	}
	if contextMessage != "" {
		content = append(content, llms.TextParts(llms.ChatMessageTypeAI, contextMessage))
	}
	resultBuilder := strings.Builder{}
	completion, err := llm.GenerateContent(llmCtx, content, llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		resultBuilder.Write(chunk)
		return nil
	}))
	cancelLLM()
	if err != nil {
		return "", err
	}
	if completion != nil {
		fmt.Printf("Complete with %d choice", len(completion.Choices))
	}
	return resultBuilder.String(), nil
}
