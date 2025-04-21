package transcription

import "github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"

type TranscribeOptions struct {
	Language         string
	ProgressCallback whisper.ProgressCallback // progress callback while processing
	SegmentCallback  whisper.SegmentCallback  // segment callback while processing
}

// CreateOptions creates a new TranscribeOptions with default values.
func CreateOptions() TranscribeOptions {
	options := TranscribeOptions{}
	options.Language = "auto"
	return options
}
