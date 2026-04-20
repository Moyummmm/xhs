package router

import (
	"server/internal/handler"
	"server/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.RouterGroup) {
	// 公开接口
	r.GET("/users/me", middleware.Auth(), handler.GetCurrentUser)
	r.GET("/users/:id", handler.GetUserById)
	r.GET("/users/:id/notes", handler.GetUserNotes)
	r.GET("/users/:id/likes", handler.GetLikedNotes)
	r.GET("/users/:id/followers", handler.GetFollowers)
	r.GET("/users/:id/followings", handler.GetFollowings)

	// 需要登录的接口
	auth := r.Group("/")
	auth.Use(middleware.Auth())
	{
		auth.PUT("/users/:id", handler.UpdateUserById)
		auth.DELETE("/users/:id", handler.DeleteUserById)
		auth.POST("/users/:id/follow", handler.FollowUser)
		auth.DELETE("/users/:id/follow", handler.UnfollowUser)
		auth.POST("/user/update", handler.Update)
	}
}
