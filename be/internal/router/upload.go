package router

import (
	"server/internal/handler"
	"server/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUploadRoutes(r *gin.RouterGroup) {
	auth := r.Group("/")
	auth.Use(middleware.Auth())
	{
		auth.POST("/upload/image", handler.UploadImage)
		auth.POST("/upload/video", handler.UploadVideo)
	}
}
