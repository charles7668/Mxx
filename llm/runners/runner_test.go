//go:build exclude_test

package runners_test

import (
	"Mxx/llm/models"
	"context"
	"testing"
)

func TestPrepareRunner(t *testing.T) {
	runner, err := PrepareRunner(Ollama, func(options *models.RunnerOptions) {
		options.ModelName = "mistral"
	})
	if err != nil {
		t.Fatalf("prepare runner error: %v", err)
	}
	ctx := context.Background()
	result, err := runner.Chat(ctx, "", "what your name", "")
	if err != nil {
		t.Fatalf("Chat() returned an error: %v", err)
	}
	t.Logf("Chat() result: %s", result)
}
