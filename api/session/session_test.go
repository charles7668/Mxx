package session

import (
	"errors"
	"testing"
	"time"
)

func TestSessionManagement(t *testing.T) {
	backupAliveTime := sessionAliveTime
	sessionAliveTime = 1 * time.Second
	defer func() {
		sessionAliveTime = backupAliveTime
	}()
	sessionId := GenerateSessionId()
	if IsAlive(sessionId) {
		t.Errorf("session should not be alive before adding")
	}
	AddToManager(sessionId, time.Now())
	if !IsAlive(sessionId) {
		t.Errorf("session should be alive")
	}
	if len(sessionMap) != 1 {
		t.Errorf("session map should have only one session")
	}
	time.Sleep(2 * time.Second)
	osRemoveAll = func(path string) error {
		return errors.New("test error")
	}
	// this case should print an error message , because osRemoveAll is mocked
	if IsAlive(sessionId) {
		t.Errorf("session should not be alive")
	}
}
