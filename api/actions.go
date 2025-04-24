package api

import (
	"Mxx/api/media"
	"Mxx/api/session"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func generateSessionId(c *gin.Context) {
	sessionId := session.GenerateSessionId()
	session.AddToManager(sessionId, time.Now())
	c.JSON(200, gin.H{"session_id": sessionId})
}

func mediaUpload(c *gin.Context) {
	sessionId := c.GetHeader("X-Session-Id")
	// create session dir if not exist
	if stat, err := os.Stat(sessionId); os.IsNotExist(err) {
		err = os.MkdirAll(sessionId, os.ModePerm)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session directory"})
			return
		}
	} else if !stat.IsDir() {
		// If the sessionId is not a directory, return an error because a required file might have the same name as this ID.
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID is invalid"})
		return
	}

	// get file from form
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}
	targetPath := filepath.Join(sessionId, file.Filename)
	if err := c.SaveUploadedFile(file, targetPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	mediaManager := media.GetMediaManager()
	mediaManager.AddMediaPath(sessionId, targetPath)
}

func generateMediaSubtitles(c *gin.Context) {
	sessionId := c.GetHeader("X-Session-Id")
	mediaManager := media.GetMediaManager()
	mediaPath := mediaManager.GetMediaPath(sessionId)
	if mediaPath == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "No media file found"})
		return
	}
	// todo : generate subtitles
}
