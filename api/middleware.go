package api

import (
	"Mxx/api/media"
	"Mxx/api/session"
	"github.com/gin-gonic/gin"
	"net/http"
)

func sessionCheckMiddleware(c *gin.Context) {
	sessionId := c.GetHeader("X-Session-Id")
	if sessionId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Session ID is required"})
		return
	}
	if !session.IsAlive(sessionId) {
		mediaManager := media.GetMediaManager()
		mediaManager.RemoveMediaPath(sessionId)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Session ID is expired"})
		return
	}
}
