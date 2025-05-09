package api

import (
	"Mxx/api/graceful"
	"Mxx/api/log"
	"embed"
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"io/fs"
	"net/http"
	"strings"
	"time"
)

var StaticFS embed.FS

func GetApiRouter(prefix string) *gin.Engine {
	graceful.InitContext()
	router := gin.New()
	prefix = strings.TrimSuffix(prefix, "/")
	apiRouterGroup := router.Group(prefix)
	// Enable CORS
	apiRouterGroup.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Replace "*" with specific origins if needed
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "X-Session-Id"},
		AllowCredentials: true,
	}))
	logger := log.GetApiLogger()
	apiRouterGroup.Use(ginzap.GinzapWithConfig(logger, &ginzap.Config{
		TimeFormat: time.RFC3339,
		UTC:        true,
		SkipPaths:  []string{"/" + prefix + "/medias/task"},
	}))
	apiRouterGroup.Use(ginzap.RecoveryWithZap(logger, true))
	apiRouterGroup.Use(prepareLogger)
	apiRouterGroup.GET("/session", generateSessionId)
	medias := apiRouterGroup.Group("/medias")
	{
		// session check middleware
		medias.Use(sessionCheckMiddleware)
		medias.Use(prepareLoggerWithSessionField)
		medias.POST("", mediaUpload)
		medias.GET("", getUploadedMedia)
		medias.POST("/subtitles", generateMediaSubtitles)
		medias.GET("/subtitles", getSubtitle)
		medias.GET("/subtitles/ass", getASSFormatSubtitle)
		medias.GET("/task", getMediaTaskState)
	}

	videos := apiRouterGroup.Group("/video/:token")
	{
		videos.Use(sessionCheckMiddleware)
		videos.Use(prepareLoggerWithSessionField)
		videos.GET("/output.m3u8", getPreviewMediaList)
		videos.GET("/:segment", getPreviewMediaFile)
	}

	return router
}

func GetWebRouter() *gin.Engine {
	router := GetApiRouter("api")

	fsys, err := fs.Sub(StaticFS, "web/dist")
	if err != nil {
		panic(err)
	}
	router.NoRoute(func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.URL.Path, "/api/") {
			fileServer := http.FileServerFS(fsys)
			c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, "/")
			fileServer.ServeHTTP(c.Writer, c.Request)
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "API not found"})
		}
	})
	return router
}
