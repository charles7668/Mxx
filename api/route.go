package api

import (
	"Mxx/api/graceful"
	"github.com/gin-gonic/gin"
)

func GetApiRouter() *gin.Engine {
	graceful.InitContext()
	router := gin.Default()
	router.GET("/session", generateSessionId)
	medias := router.Group("/medias")
	{
		// session check middleware
		medias.Use(sessionCheckMiddleware)
		medias.POST("", mediaUpload)
		medias.POST("/subtitles", generateMediaSubtitles)
		medias.GET("/subtitles", getSubtitle)
		medias.GET("/task", getMediaTaskState)
	}
	return router
}
