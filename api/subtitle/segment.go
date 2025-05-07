package subtitle

import "time"

type Segment struct {
	StartTime time.Duration `json:"start_time"`
	EndTime   time.Duration `json:"end_time"`
	Text      string        `json:"text"`
}
