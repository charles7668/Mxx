package models

import "context"

type ChatRunner interface {
	Chat(ctx context.Context, characterMessage, prompt, contextMessage string) (string, error)
}
