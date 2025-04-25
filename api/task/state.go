package task

import "context"

type RunningStatus int

const (
	Running RunningStatus = iota
	Completed
	Failed
)

type State struct {
	Status     RunningStatus
	Task       string
	CancelFunc context.CancelFunc
}

func (s *State) String() string {
	switch s.Status {
	case Running:
		return "Running"
	case Completed:
		return "Completed"
	case Failed:
		return "Failed"
	default:
		return "Completed"
	}
}
