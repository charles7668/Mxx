package api

import (
	"Mxx/api/graceful"
	"Mxx/api/log"
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"time"
)

func GetApiRouter() *gin.Engine {
	graceful.InitContext()
	router := gin.New()
	// Enable CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Replace "*" with specific origins if needed
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "X-Session-Id"},
		AllowCredentials: true,
	}))
	logger := log.GetLogger()
	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(logger, true))
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
	router.GET("/video/:token/output.m3u8", getPreviewMediaList)
	router.GET("/video/:token/:segment", getPreviewMediaFile)

	return router
}
