package session

import (
	"fmt"
	"os"
	"time"
)

var sessionAliveTime = 10 * time.Minute

// this map is used to store the session id
var sessionMap = make(map[string]time.Time)

func Update(sessionId string, timestamp time.Time) {
	sessionMap[sessionId] = timestamp
}

// os.RemoveAll function , extract here for cover os.RemoveAll fail case
var osRemoveAll = os.RemoveAll

func IsAlive(sessionId string) bool {
	_, ok := sessionMap[sessionId]
	if !ok {
		return false
	}
	// check if the session is alive
	if time.Since(sessionMap[sessionId]) > sessionAliveTime {
		delete(sessionMap, sessionId)
		err := osRemoveAll(sessionId)
		if err != nil {
			fmt.Printf("failed to remove directory for session : %s , err : %s \n", sessionId, err.Error())
		}
		return false
	}
	return true
}
