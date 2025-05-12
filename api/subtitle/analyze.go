package subtitle

import (
	"strings"
	"time"
)

func TryMerge(previous, current *Segment, language string) (*Segment, bool) {
	if language == "zh" {
		return nil, false
	}
	result := &Segment{
		StartTime: previous.StartTime,
		EndTime:   current.EndTime,
		Text:      previous.Text,
	}
	trim := strings.TrimSpace(current.Text)
	if len(trim) == 0 {
		return nil, false
	}
	if isEndOfSentence(strings.TrimSpace(previous.Text)) {
		return nil, false
	}
	if isEndOfSentence(trim) {
		result.EndTime = current.StartTime
		result.Text = strings.TrimSuffix(previous.Text, " ") + trim
		return result, true
	}
	if isPausePoint(trim) {
		result.EndTime = current.EndTime
		result.Text = strings.TrimSuffix(previous.Text, " ") + trim
		return result, true
	}
	if IsBrace(trim) {
		result.EndTime = current.EndTime
		result.Text = strings.TrimSuffix(previous.Text, " ") + " " + trim + " "
		return result, true
	}
	if current.StartTime-previous.EndTime < 100*time.Millisecond {
		result.EndTime = current.EndTime
		if trim[0] == '\'' {
			result.Text = previous.Text + trim
		} else {
			result.Text = previous.Text + " " + trim
		}
		return result, true
	}
	return nil, false
}

func IsBrace(text string) bool {
	if len(text) != 1 {
		return false
	}
	for _, r := range text {
		if !strings.ContainsRune("()[]{}<>", r) {
			return false
		}
	}
	return true
}

func isPausePoint(text string) bool {
	if len(text) != 1 {
		return false
	}
	for _, r := range text {
		if !strings.ContainsRune(",:", r) {
			return false
		}
	}
	return true
}

func isEndOfSentence(text string) bool {
	for _, r := range text[len(text)-1:] {
		if !strings.ContainsRune(".!?;", r) {
			return false
		}
	}
	return true
}
