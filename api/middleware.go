package api

import (
	"Mxx/api/constant"
	"Mxx/api/media"
	"Mxx/api/session"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func sessionCheckMiddleware(c *gin.Context) {
	sessionId := c.GetHeader("X-Session-Id")
	if sessionId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Session ID is required"})
		c.Abort()
		return
	}
	if !session.IsAlive(sessionId) {
		mediaManager := media.GetMediaManager()
		mediaManager.RemoveMediaPath(sessionId)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Session ID is expired"})
		c.Abort()
		return
	}
	c.Set(constant.SessionIdCtxKey, sessionId)
	// renew session
	session.Update(sessionId, time.Now())
}
