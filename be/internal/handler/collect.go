package handler

import (
	"strconv"
	"server/config"
	"server/internal/middleware"
	"server/internal/repository"
	"server/internal/service"
	"server/pkg/errorConfig"
	"server/pkg/response"

	"github.com/gin-gonic/gin"
)

var collectService *service.CollectService

func init() {
	collectRepo := repository.NewCollectRepository(config.DB())
	collectService = service.NewCollectService(collectRepo)
}

// Post /notes/:id/collect
// AddCollect 用户收藏 Note
func AddCollect(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		response.Fail(c, errorConfig.ErrBadRequest.Code, "FE参数错误")
		return
	}
	noteId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Fail(c, errorConfig.ErrBadRequest.Code, "FE参数错误")
		return
	}
	// 获取userId
	userId := uint(middleware.CurrentUserID(c))
	if userId == 0 {
		response.Fail(c, errorConfig.ErrUnauthorized.Code, "未登录")
		return
	}
	err = collectService.CollectById(userId, uint(noteId))
	if err != nil {
		response.Fail(c, errorConfig.ErrInternalServer.Code, "收藏失败")
		return
	}
	response.Success(c, "收藏成功")
}

// Get /note/collect
// GetCollectList 获取用户收藏列表
func GetCollectList(c *gin.Context) {
	userId := uint(middleware.CurrentUserID(c))
	if userId == 0 {
		response.Fail(c, errorConfig.ErrUnauthorized.Code, "未登录")
		return
	}

	collectList, err := collectService.GetCollectList(userId)
	if err != nil {
		response.Fail(c, errorConfig.ErrInternalServer.Code, "获取收藏列表失败")
		return
	}
	response.Success(c, collectList)
}	

// Get /note/collect/count
// GetCollectedCount 获取用户收藏数量
func GetCollectedCount(c *gin.Context) {
	userId := uint(middleware.CurrentUserID(c))
	if userId == 0 {
		response.Fail(c, errorConfig.ErrUnauthorized.Code, "未登录")
		return
	}

	count, err := collectService.GetCollectedCount(userId)
	if err != nil {
		response.Fail(c, errorConfig.ErrInternalServer.Code, "获取收藏数量失败")
		return
	}
	response.Success(c, count)
}

// Delete /note/collect/:id
// DisCollectById 用户取消收藏 Note
func DisCollectById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Fail(c, errorConfig.ErrBadRequest.Code, "FE参数错误")
		return
	}
	noteId, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.Fail(c, errorConfig.ErrBadRequest.Code, "FE参数错误")
		return
	}
	// 获取userId
	userId := uint(middleware.CurrentUserID(c))
	if userId == 0 {
		response.Fail(c, errorConfig.ErrUnauthorized.Code, "未登录")
		return
	}
	err = collectService.DisCollectById(userId, uint(noteId))
	if err != nil {
		response.Fail(c, errorConfig.ErrInternalServer.Code, "取消收藏失败")
		return
	}
	response.Success(c, "取消收藏成功")
}