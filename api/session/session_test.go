package session

import (
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
	AddToManager(sessionId, time.Now())
	if !IsAlive(sessionId) {
		t.Errorf("session should be alive")
	}
	if len(sessionMap) != 1 {
		t.Errorf("session map should have only one session")
	}
	time.Sleep(2 * time.Second)
	if IsAlive(sessionId) {
		t.Errorf("session should not be alive")
	}
}
