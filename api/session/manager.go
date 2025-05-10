package session

import (
	"fmt"
	"os"
	"sync"
	"time"
)

var sessionAliveTime = 10 * time.Minute

var mu = &sync.Mutex{}

// this map is used to store the session id
var sessionMap = make(map[string]time.Time)

func Update(sessionId string, timestamp time.Time) {
	mu.Lock()
	defer mu.Unlock()
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
		mu.Lock()
		defer mu.Unlock()
		delete(sessionMap, sessionId)
		err := osRemoveAll(sessionId)
		if err != nil {
			fmt.Printf("failed to remove directory for session : %s , err : %s \n", sessionId, err.Error())
		}
		return false
	}
	return true
}
