package api

import (
	"github.com/gin-gonic/gin"
)

func GetApiRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/session", generateSessionId)
	medias := router.Group("/medias")
	{
		// session check middleware
		medias.Use(sessionCheckMiddleware)
		medias.POST("", mediaUpload)
		medias.GET("/subtitles", generateMediaSubtitles)
	}
	return router
}
