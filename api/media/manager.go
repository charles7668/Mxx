package media

type Manager struct {
	mediaRecords map[string]string
}

var mediaManager *Manager

// AddMediaPath Add media path to the map for management
func (m *Manager) AddMediaPath(sessionId, path string) {
	m.mediaRecords[sessionId] = path
}

// RemoveMediaPath Remove media path from the map
func (m *Manager) RemoveMediaPath(sessionId string) {
	delete(m.mediaRecords, sessionId)
}

func (m *Manager) GetMediaPath(sessionId string) string {
	if path, ok := m.mediaRecords[sessionId]; ok {
		return path
	}
	return ""
}

func GetMediaManager() *Manager {
	if mediaManager == nil {
		mediaManager = &Manager{
			mediaRecords: make(map[string]string),
		}
	}
	return mediaManager
}
