package api

import (
	"Mxx/api/session"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func GetApiRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/session", func(c *gin.Context) {
		sessionId := session.GenerateSessionId()
		session.AddToManager(sessionId, time.Now())
		c.JSON(200, gin.H{"session_id": sessionId})
	})
	medias := router.Group("/medias")
	{
		medias.POST("", func(c *gin.Context) {
			sessionId := c.GetHeader("X-Session-Id")
			if sessionId == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Session ID is required"})
				return
			}
			if !session.IsAlive(sessionId) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Session ID is expired"})
				return
			}
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
		})
	}
	return router
}
