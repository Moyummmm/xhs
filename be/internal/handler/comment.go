package handler

import (
	"net/http"
	"server/config"
	"server/internal/middleware"
	"server/internal/repository"
	"server/internal/service"
	"server/pkg/errorConfig"
	"server/pkg/jwt"
	"server/pkg/response"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var commentService *service.CommentService

func init() {
	commentRepo := repository.NewCommentRepository(config.DB())
	commentService = service.NewCommentService(commentRepo)
}

func extractOptionalUserID(c *gin.Context) uint {
	tokenString := extractTokenFromRequest(c)
	if tokenString == "" {
		return 0
	}
	claims, err := jwt.ParseToken(tokenString)
	if err != nil {
		return 0
	}
	return uint(claims.UserID)
}

func extractTokenFromRequest(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
			return parts[1]
		}
		return authHeader
	}

	if token := c.GetHeader("X-Token"); token != "" {
		return token
	}

	if token := c.GetHeader("AccessToken"); token != "" {
		return token
	}

	return c.Query("token")
}

func GetComments(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	noteID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Fail(c, errorConfig.ErrBadRequest.Code, "笔记ID格式错误")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	sort := c.DefaultQuery("type", "")
	if sort == "" {
		sort = c.DefaultQuery("sort", "hot")
	}
	if sort != "hot" && sort != "latest" {
		sort = "hot"
	}

	userID := extractOptionalUserID(c)

	result, err := commentService.GetComments(ctx, uint(noteID), page, pageSize, sort, userID)
	if err != nil {
		response.Fail(c, errorConfig.ErrInternalServer.Code, "获取评论失败")
		return
	}

	response.Success(c, result)
}

type CreateCommentReq struct {
	Content  string `json:"content" binding:"required,max=500"`
	ParentID *uint  `json:"parent_id"`
}

func CreateComment(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	noteID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Fail(c, errorConfig.ErrBadRequest.Code, "笔记ID格式错误")
		return
	}

	var req CreateCommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errorConfig.ErrBadRequest.Code, "请求参数错误")
		return
	}

	userID := uint(middleware.CurrentUserID(c))
	if userID == 0 {
		response.Fail(c, errorConfig.ErrUnauthorized.Code, "未登录")
		return
	}

	if err := commentService.CreateComment(ctx, userID, uint(noteID), req.Content, req.ParentID); err != nil {
		if err == service.ErrInvalidContent {
			response.Fail(c, errorConfig.ErrBadRequest.Code, "评论内容不能为空")
			return
		}
		response.Fail(c, errorConfig.ErrInternalServer.Code, "发布评论失败")
		return
	}

	response.Success(c, "发布成功")
}

func DeleteComment(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	noteID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Fail(c, errorConfig.ErrBadRequest.Code, "笔记ID格式错误")
		return
	}
	_ = noteID

	commentIDStr := c.Param("comment_id")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 64)
	if err != nil {
		response.Fail(c, errorConfig.ErrBadRequest.Code, "评论ID格式错误")
		return
	}

	userID := uint(middleware.CurrentUserID(c))
	if userID == 0 {
		response.Fail(c, errorConfig.ErrUnauthorized.Code, "未登录")
		return
	}

	if err := commentService.DeleteComment(ctx, userID, uint(commentID)); err != nil {
		if err == service.ErrCommentNotFound {
			response.FailWithStatus(c, http.StatusNotFound, errorConfig.ErrNotFound.Code, "评论不存在")
			return
		}
		if err == service.ErrNoPermission {
			response.FailWithStatus(c, http.StatusForbidden, errorConfig.ErrForbidden.Code, "无权限删除")
			return
		}
		response.Fail(c, errorConfig.ErrInternalServer.Code, "删除评论失败")
		return
	}

	response.Success(c, "删除成功")
}
