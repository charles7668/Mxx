package task

import "context"

const (
	Running = iota
	Completed
	Failed
)

type State struct {
	RunningStatus int
	State         string
	CancelFunc    context.CancelFunc
}
