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
	followService = service.NewFollowService(followRepo, userRepo)
}

type UpdateRequest model.User

func Update(c *gin.Context) {
	ctx := c.Request.Context()

	var req UpdateRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		response.Fail(c, errorConfig.ErrBadRequest.Code, "request parameter wrong")
		return
	}
	result, _ := userService.Patch(ctx, model.User(req))
	if result {
		response.Success(c, "success")
		return
	}
	response.Fail(c, 500, "内部错误")
}

func GetCurrentUser(c *gin.Context) {
	ctx := c.Request.Context()

	userId := uint(middleware.CurrentUserID(c))
	if userId == 0 {
		response.Fail(c, errorConfig.ErrUnauthorized.Code, "未登录")
		return
	}

	user, err := userService.GetById(ctx, userId)
	if err != nil {
		response.Fail(c, errorConfig.ErrInternalServer.Code, "获取用户信息失败")
		return
	}
	response.Success(c, user)
}

func GetUserById(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Fail(c, errorConfig.ErrBadRequest.Code, "用户ID格式错误")
		return
	}

	user, err := userService.GetById(ctx, uint(id))
	if err != nil {
		response.Fail(c, errorConfig.ErrInternalServer.Code, "获取用户信息失败")
		return
	}
	response.Success(c, user)
}

func UpdateUserById(c *gin.Context) {
	ctx := c.Request.Context()

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

	updated, err := userService.UpdateById(ctx, uint(id), req)
	if err != nil {
		response.Fail(c, errorConfig.ErrInternalServer.Code, "更新用户信息失败")
		return
	}
	response.Success(c, updated)
}

func DeleteUserById(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Fail(c, errorConfig.ErrBadRequest.Code, "用户ID格式错误")
		return
	}

	if err := userService.DeleteById(ctx, uint(id)); err != nil {
		response.Fail(c, errorConfig.ErrInternalServer.Code, "删除用户失败")
		return
	}
	response.Success(c, "删除成功")
}

func FollowUser(c *gin.Context) {
	ctx := c.Request.Context()

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
	if err := followService.Follow(ctx, userId, uint(targetId)); err != nil {
		response.Fail(c, errorConfig.ErrInternalServer.Code, "关注失败")
		return
	}
	response.Success(c, "关注成功")
}

func UnfollowUser(c *gin.Context) {
	ctx := c.Request.Context()

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
	if err := followService.Unfollow(ctx, userId, uint(targetId)); err != nil {
		response.Fail(c, errorConfig.ErrInternalServer.Code, "取消关注失败")
		return
	}
	response.Success(c, "取消关注成功")
}

func GetFollowers(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	userId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Fail(c, errorConfig.ErrBadRequest.Code, "用户ID格式错误")
		return
	}

	users, err := followService.GetFollowers(ctx, uint(userId))
	if err != nil {
		response.Fail(c, errorConfig.ErrInternalServer.Code, "获取粉丝列表失败")
		return
	}
	response.Success(c, users)
}

func GetFollowings(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	userId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Fail(c, errorConfig.ErrBadRequest.Code, "用户ID格式错误")
		return
	}

	users, err := followService.GetFollowings(ctx, uint(userId))
	if err != nil {
		response.Fail(c, errorConfig.ErrInternalServer.Code, "获取关注列表失败")
		return
	}
	response.Success(c, users)
}
