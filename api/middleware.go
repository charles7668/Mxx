package api

import (
	"Mxx/api/constant"
	"Mxx/api/log"
	"Mxx/api/media"
	"Mxx/api/session"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func sessionCheckMiddleware(c *gin.Context) {
	sessionId := c.GetHeader("X-Session-Id")
	if sessionId == "" {
		// if no session id in header, check if token is in url
		token := c.Param("token")
		sessionId = token
		if sessionId == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Session ID is required"})
			c.Abort()
			return
		}
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

func prepareLogger(c *gin.Context) {
	logger := log.GetInnerLogger()
	c.Set(constant.LoggerCtxKey, logger)
}

func prepareLoggerWithSessionField(c *gin.Context) {
	sessionId := c.GetString(constant.SessionIdCtxKey)
	if sessionId == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	logger := log.GetInnerLogger()
	newLogger := logger.With(zap.String("sessionId", sessionId))
	c.Set(constant.LoggerCtxKey, newLogger)
}
