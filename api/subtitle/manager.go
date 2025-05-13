package subtitle

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type Manager struct {
	cache map[string][]Segment
	mu    *sync.Mutex
}

var (
	manager *Manager
)

func GetManager() *Manager {
	if manager == nil {
		manager = &Manager{
			cache: make(map[string][]Segment),
			mu:    &sync.Mutex{},
		}
	}
	return manager
}

func (m *Manager) Add(sessionId string, segment Segment) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.cache[sessionId]; !ok {
		m.cache[sessionId] = []Segment{}
	}
	m.cache[sessionId] = append(m.cache[sessionId], segment)
}

func (m *Manager) Clear(sessionId string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.cache, sessionId)
}

func (m *Manager) GetSegments(sessionId string) []Segment {
	segments, ok := m.cache[sessionId]
	if !ok {
		return []Segment{}
	}
	return segments
}

func (m *Manager) Exist(sessionId string) bool {
	segments, ok := m.cache[sessionId]
	if !ok || len(segments) == 0 {
		return false
	}
	return true
}

func duratiionToASSTimeFormat(d *time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	totalCentiseconds := d.Milliseconds() / 10
	centiSeconds := int(totalCentiseconds) % 100

	return fmt.Sprintf("%02d:%02d:%02d.%02d", hours, minutes, seconds, centiSeconds)
}

func (m *Manager) ToASS(sessionId string) string {
	var builder strings.Builder

	builder.WriteString("[Script Info]\n")
	builder.WriteString("Title: Subtitles\n")
	builder.WriteString("Original Script: Mxx\n")
	builder.WriteString("ScriptType: v4.00+\n")
	builder.WriteString("PlayResX: 1280\n")
	builder.WriteString("PlayResY: 720\n")
	builder.WriteString("Collisions: Normal\n")
	builder.WriteString("WrapStyle: 0\n")
	builder.WriteString("ScaledBorderAndShadow: yes\n")

	builder.WriteString("[V4+ Styles]\n")
	builder.WriteString("Format: Name, Fontname, Fontsize, PrimaryColour, SecondaryColour, OutlineColour, BackColour, Bold, Italic, Underline, StrikeOut, ScaleX, ScaleY, Spacing, Angle, BorderStyle, Outline, Shadow, Alignment, MarginL, MarginR, MarginV, Encoding\n")
	builder.WriteString("Style: Default,Arial,36,&H00FFFFFF,&H000000FF,&H00000000,&H64000000,-1,0,0,0,100,100,0,0,1,2,2,2,10,10,10,1\n")

	builder.WriteString("[Events]\n")
	builder.WriteString("Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text\n")
	for _, segment := range m.cache[sessionId] {
		builder.WriteString(fmt.Sprintf("Dialogue: 0,%s,%s,Default,,0,0,0,,%s\n", duratiionToASSTimeFormat(&segment.StartTime), duratiionToASSTimeFormat(&segment.EndTime), segment.Text))
	}
	result := builder.String()
	return result
}

func (m *Manager) Last(sessionId string) *Segment {
	m.mu.Lock()
	defer m.mu.Unlock()
	segments, ok := m.cache[sessionId]
	if !ok || len(segments) == 0 {
		return nil
	}
	return &segments[len(segments)-1]
}

func (m *Manager) ToPlainText(sessionId string) string {
	builder := strings.Builder{}
	for _, segment := range m.GetSegments(sessionId) {
		builder.WriteString("[")
		builder.WriteString(segment.StartTime.String())
		builder.WriteString(" -> ")
		builder.WriteString(segment.EndTime.String())
		builder.WriteString("] ")
		builder.WriteString(segment.Text)
		builder.WriteString("\n")
	}

	return builder.String()
}
