package session

import "github.com/google/uuid"

func GenerateSessionId() string {
	id := uuid.New()
	return id.String()
}
