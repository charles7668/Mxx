package subtitle

import "sync"

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
	_, ok := m.cache[sessionId]
	return ok
}
