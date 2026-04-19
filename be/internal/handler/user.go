package handler

import (
	"server/config"
	"server/internal/middleware"
	"server/internal/model"
	"server/internal/repository"
	"server/internal/service"
	"server/pkg/errorConfig"
	"server/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

var userService *service.UserService
var followService *service.FollowService

func init() {
	userRepo := repository.NewUserRepository(config.DB())
	userService = service.NewUserService(userRepo)

	followRepo := repository.NewFollowRepository(config.DB())
	followService = service.NewFollowService(followRepo)
}

type UpdateRequest model.User

// Method: POST /user/update
// Update 用户修改信息（保留以兼容旧接口）
func Update(c *gin.Context) {
	var req UpdateRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		response.Fail(c, errorConfig.ErrBadRequest.Code, "request parameter wrong")
		return
	}
	result, _ := userService.Patch(model.User(req))
	if result {
		response.Success(c, "success")
		return
	}
	response.Fail(c, 500, "内部错误")
}

// GetCurrentUser 获取当前登录用户信息
// GET /users/me
func GetCurrentUser(c *gin.Context) {
	userId := uint(middleware.CurrentUserID(c))
	if userId == 0 {
		response.Fail(c, errorConfig.ErrUnauthorized.Code, "未登录")
		return
	}

	user, err := userService.GetById(userId)
	if err != nil {
		response.Fail(c, errorConfig.ErrInternalServer.Code, "获取用户信息失败")
		return
	}
	response.Success(c, user)
}

// GetUserById 获取用户信息
// GET /users/:id
func GetUserById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Fail(c, errorConfig.ErrBadRequest.Code, "用户ID格式错误")
		return
	}

	user, err := userService.GetById(uint(id))
	if err != nil {
		response.Fail(c, errorConfig.ErrInternalServer.Code, "获取用户信息失败")
		return
	}
	response.Success(c, user)
}

// UpdateUserById 更新用户信息
// PUT /users/:id
func UpdateUserById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Fail(c, errorConfig.ErrBadRequest.Code, "用户ID格式错误")
		return
	}

	var req model.User
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errorConfig.ErrBadRequest.Code, "请求参数错误")
		return
	}

	if err := userService.UpdateById(uint(id), req); err != nil {
		response.Fail(c, errorConfig.ErrInternalServer.Code, "更新用户信息失败")
		return
	}
	response.Success(c, "更新成功")
}

// DeleteUserById 删除用户
// DELETE /users/:id
func DeleteUserById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Fail(c, errorConfig.ErrBadRequest.Code, "用户ID格式错误")
		return
	}

	if err := userService.DeleteById(uint(id)); err != nil {
		response.Fail(c, errorConfig.ErrInternalServer.Code, "删除用户失败")
		return
	}
	response.Success(c, "删除成功")
}

// FollowUser 关注用户
// POST /users/:id/follow
func FollowUser(c *gin.Context) {
	idStr := c.Param("id")
	targetId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Fail(c, errorConfig.ErrBadRequest.Code, "用户ID格式错误")
		return
	}

	userId := uint(middleware.CurrentUserID(c))
	if userId == 0 {
		response.Fail(c, errorConfig.ErrUnauthorized.Code, "未登录")
		return
	}
	if err := followService.Follow(userId, uint(targetId)); err != nil {
		response.Fail(c, errorConfig.ErrInternalServer.Code, "关注失败")
		return
	}
	response.Success(c, "关注成功")
}

// UnfollowUser 取消关注用户
// DELETE /users/:id/follow
func UnfollowUser(c *gin.Context) {
	idStr := c.Param("id")
	targetId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Fail(c, errorConfig.ErrBadRequest.Code, "用户ID格式错误")
		return
	}

	userId := uint(middleware.CurrentUserID(c))
	if userId == 0 {
		response.Fail(c, errorConfig.ErrUnauthorized.Code, "未登录")
		return
	}
	if err := followService.Unfollow(userId, uint(targetId)); err != nil {
		response.Fail(c, errorConfig.ErrInternalServer.Code, "取消关注失败")
		return
	}
	response.Success(c, "取消关注成功")
}

// GetFollowers 获取粉丝列表
// GET /users/:id/followers
func GetFollowers(c *gin.Context) {
	idStr := c.Param("id")
	userId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Fail(c, errorConfig.ErrBadRequest.Code, "用户ID格式错误")
		return
	}

	users, err := followService.GetFollowers(uint(userId))
	if err != nil {
		response.Fail(c, errorConfig.ErrInternalServer.Code, "获取粉丝列表失败")
		return
	}
	response.Success(c, users)
}

// GetFollowings 获取关注列表
// GET /users/:id/followings
func GetFollowings(c *gin.Context) {
	idStr := c.Param("id")
	userId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Fail(c, errorConfig.ErrBadRequest.Code, "用户ID格式错误")
		return
	}

	users, err := followService.GetFollowings(uint(userId))
	if err != nil {
		response.Fail(c, errorConfig.ErrInternalServer.Code, "获取关注列表失败")
		return
	}
	response.Success(c, users)
}
