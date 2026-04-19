package router

import (
	"server/internal/handler"
	"server/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterNoteRoutes(r *gin.RouterGroup) {
	// 公开接口
	r.GET("/notes/feed", handler.GetNote)
	r.GET("/notes/:id", handler.GetNoteById)
	r.GET("/notes/search", handler.SearchNotes)
	r.GET("/notes/:id/comments", handler.GetComments)

	// 需要登录的接口
	auth := r.Group("/")
	auth.Use(middleware.Auth())
	{
		auth.POST("/notes", handler.CreateNote)
		auth.PUT("/notes/:id", handler.UpdateNote)
		auth.DELETE("/notes/:id", handler.DeleteByNoteId)
		auth.POST("/notes/:id/like", handler.LikeNote)
		auth.DELETE("/notes/:id/like", handler.UnlikeNote)
		auth.POST("/notes/:id/comments", handler.CreateComment)
		auth.DELETE("/notes/:id/comments/:comment_id", handler.DeleteComment)
	}
}
