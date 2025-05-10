package models

import "time"

type Segment struct {
	Start    time.Duration
	End      time.Duration
	Duration time.Duration
}
