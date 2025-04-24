package task

import (
	"github.com/patrickmn/go-cache"
	"sync"
)

var taskState = cache.New(cache.NoExpiration, cache.NoExpiration)
var completeTask = cache.New(cache.NoExpiration, cache.NoExpiration)
var failedTask = cache.New(cache.NoExpiration, cache.NoExpiration)
var mutex sync.Mutex

func StartTask(taskId string, state State) {
	state.RunningStatus = Running
	taskState.Set(taskId, state, cache.NoExpiration)
	mutex.Lock()
	defer mutex.Unlock()
	failedTask.Delete(taskId)
	completeTask.Delete(taskId)
}

func FailedTask(taskId string, error error) {
	mutex.Lock()
	defer mutex.Unlock()
	state, found := GetTaskState(taskId)
	if !found {
		return
	}
	if state.CancelFunc != nil {
		state.CancelFunc()
	}
	state.State = error.Error()
	state.RunningStatus = Failed
	failedTask.Set(taskId, state, cache.NoExpiration)
	taskState.Delete(taskId)
}

func CompleteTask(taskId string) {
	mutex.Lock()
	defer mutex.Unlock()
	state, found := GetTaskState(taskId)
	if !found {
		return
	}
	if state.CancelFunc != nil {
		state.CancelFunc()
	}
	state.RunningStatus = Completed
	completeTask.Set(taskId, state, cache.NoExpiration)
	taskState.Delete(taskId)
}

func GetTaskState(taskId string) (state State, found bool) {
	value, found := taskState.Get(taskId)
	if found {
		return value.(State), true
	}
	value, found = completeTask.Get(taskId)
	if found {
		return value.(State), true
	}
	value, found = failedTask.Get(taskId)
	if found {
		return value.(State), true
	}
	return State{}, false
}
