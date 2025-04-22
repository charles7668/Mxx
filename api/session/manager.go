package session

import (
	"fmt"
	"os"
	"time"
)

var sessionAliveTime = 10 * time.Minute

// this map is used to store the session id
var sessionMap = make(map[string]time.Time)

func AddToManager(sessionId string, timestamp time.Time) {
	sessionMap[sessionId] = timestamp
}

func IsAlive(sessionId string) bool {
	_, ok := sessionMap[sessionId]
	if !ok {
		return false
	}
	// check if the session is alive
	if time.Since(sessionMap[sessionId]) > sessionAliveTime {
		delete(sessionMap, sessionId)
		err := os.RemoveAll(sessionId)
		if err != nil {
			fmt.Printf("failed to remove directory for session : %s , err : %s ", sessionId, err.Error())
		}
		return false
	}
	return true
}
