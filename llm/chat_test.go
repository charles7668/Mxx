//go:build exclude_test

package llm

import (
	"context"
	"testing"
)

func TestPrepareRunner(t *testing.T) {
	runner, err := PrepareRunner(Ollama, map[string]string{
		"model": "mistral",
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
