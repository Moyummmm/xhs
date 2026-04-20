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

func AddCollect(c *gin.Context) {
	ctx := c.Request.Context()

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
	userId := uint(middleware.CurrentUserID(c))
	if userId == 0 {
		response.Fail(c, errorConfig.ErrUnauthorized.Code, "未登录")
		return
	}
	err = collectService.CollectById(ctx, userId, uint(noteId))
	if err != nil {
		response.Fail(c, errorConfig.ErrInternalServer.Code, "收藏失败")
		return
	}
	response.Success(c, "收藏成功")
}

func GetCollectList(c *gin.Context) {
	ctx := c.Request.Context()

	userId := uint(middleware.CurrentUserID(c))
	if userId == 0 {
		response.Fail(c, errorConfig.ErrUnauthorized.Code, "未登录")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	collectList, total, err := collectService.GetCollectListWithPagination(ctx, userId, page, pageSize)
	if err != nil {
		response.Fail(c, errorConfig.ErrInternalServer.Code, "获取收藏列表失败")
		return
	}

	totalPage := total / int64(pageSize)
	if total%int64(pageSize) != 0 {
		totalPage++
	}

	response.Success(c, map[string]interface{}{
		"list": collectList,
		"pagination": map[string]interface{}{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
			"has_more":  page < int(totalPage),
		},
	})
}

func GetCollectedCount(c *gin.Context) {
	ctx := c.Request.Context()

	userId := uint(middleware.CurrentUserID(c))
	if userId == 0 {
		response.Fail(c, errorConfig.ErrUnauthorized.Code, "未登录")
		return
	}

	count, err := collectService.GetCollectedCount(ctx, userId)
	if err != nil {
		response.Fail(c, errorConfig.ErrInternalServer.Code, "获取收藏数量失败")
		return
	}
	response.Success(c, count)
}

func DisCollectById(c *gin.Context) {
	ctx := c.Request.Context()

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
	userId := uint(middleware.CurrentUserID(c))
	if userId == 0 {
		response.Fail(c, errorConfig.ErrUnauthorized.Code, "未登录")
		return
	}
	err = collectService.DisCollectById(ctx, userId, uint(noteId))
	if err != nil {
		response.Fail(c, errorConfig.ErrInternalServer.Code, "取消收藏失败")
		return
	}
	response.Success(c, "取消收藏成功")
}
