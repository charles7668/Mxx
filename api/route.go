package api

import (
	"Mxx/api/graceful"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func GetApiRouter() *gin.Engine {
	graceful.InitContext()
	router := gin.Default()
	// Enable CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Replace "*" with specific origins if needed
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "X-Session-Id"},
		AllowCredentials: true,
	}))
	router.GET("/session", generateSessionId)
	medias := router.Group("/medias")
	{
		// session check middleware
		medias.Use(sessionCheckMiddleware)
		medias.POST("", mediaUpload)
		medias.GET("", getUploadedMedia)
		medias.POST("/subtitles", generateMediaSubtitles)
		medias.GET("/subtitles", getSubtitle)
		medias.GET("/task", getMediaTaskState)
	}
	return router
}
