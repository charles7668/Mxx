package media

import (
	"Mxx/api/configs"
	"fmt"
	"os"
	"strings"
)

type Manager struct {
	mediaRecords map[string]string
}

var mediaManager *Manager

// SetMediaPath Set a media path to the map for management
func (m *Manager) SetMediaPath(sessionId, path string) {
	// try to remove old media in filesystem
	if oldPath, ok := m.mediaRecords[sessionId]; ok {
		// check if the old path starts with a media store path, if not then prevent deleting file
		apiConfig := configs.GetApiConfig()
		if strings.HasPrefix(oldPath, apiConfig.MediaStorePath) {
			go func() {
				err := os.Remove(oldPath)
				if err != nil {
					err = fmt.Errorf("failed to remove old media file: %s, err: %s", oldPath, err)
					fmt.Println(err)
				}
			}()
		}
	}
	// add a new media path to the map
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
