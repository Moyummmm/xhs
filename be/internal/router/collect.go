package router

import (
	"server/internal/handler"
	"server/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterCollectRoutes(r *gin.RouterGroup) {
	// 需要登录的接口
	auth := r.Group("/")
	auth.Use(middleware.Auth())
	{
		auth.POST("/notes/:id/collect", handler.AddCollect)
		auth.DELETE("/notes/:id/collect", handler.DisCollectById)
		auth.GET("/collects", handler.GetCollectList)
		auth.GET("/collects/count", handler.GetCollectedCount)
	}
}
