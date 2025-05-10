package models

import "time"

type SilentSegment Segment

// SilentAnalyzeOptions represents the options for silent analysis. if you don't know how to set it, just use the GetDefaultSilentAnalyzeOptions() function to get.
type SilentAnalyzeOptions struct {
	InputFilePath string
	NoiseDB       float64       // exp: -30.0
	Duration      time.Duration // exp: 100 * time.Millisecond
}

func GetDefaultSilentAnalyzeOptions(inputFilePath string) SilentAnalyzeOptions {
	return SilentAnalyzeOptions{
		InputFilePath: inputFilePath,
		NoiseDB:       -30.0,
		Duration:      100 * time.Millisecond,
	}
}
