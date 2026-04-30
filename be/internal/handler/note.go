package handler

import (
	"log"
	"server/config"
	"server/internal/cache"
	"server/internal/middleware"
	"server/internal/model"
	"server/internal/repository"
	"server/internal/service"
	"server/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

var noteService *service.NoteService
var likeService *service.LikeService

func init() {
	noteRepository := repository.NewNoteRepository(config.DB())
	noteService = service.NewNoteService(noteRepository)

	likeRepo := repository.NewLikeRepository(config.DB())
	likeService = service.NewLikeService(likeRepo)
}

type CreateNoteReq struct {
	Title    string `json:"title" binding:"required,max=100"`
	Content  string `json:"content" binding:"required,max=500"`
	ImageIDs []uint `json:"image_ids"`
	VideoURL string `json:"video_url"`
	Location string `json:"location"`
	TopicID  uint   `json:"topic_id"`
}

type CreateNoteResp struct {
	ID uint `json:"id"`
}

func CreateNote(c *gin.Context) {
	ctx := c.Request.Context()

	var req CreateNoteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "请求参数错误")
		return
	}

	userID := uint(middleware.CurrentUserID(c))
	if userID == 0 {
		response.Fail(c, 401, "未登录")
		return
	}

	note := &model.Note{
		Title:    req.Title,
		Body:     req.Content,
		VideoURL: req.VideoURL,
		Location: req.Location,
		TopicID:  req.TopicID,
		UserID:   userID,
	}

	if len(req.ImageIDs) > 0 {
		images, err := imageService.GetByIds(req.ImageIDs)
		if err == nil {
			for _, img := range images {
				note.Images = append(note.Images, model.NoteImage{
					NoteID: note.ID,
					URL:    img.URL,
					Width:  img.Width,
					Height: img.Height,
				})
			}
		}
	}

	if err := noteService.CreateWithImages(ctx, note); err != nil {
		response.Fail(c, 500, "创建笔记失败")
		return
	}

	response.Success(c, CreateNoteResp{ID: note.ID})
}

type UpdateNoteReq struct {
	Title    string `json:"title" binding:"required,max=100"`
	Content  string `json:"content" binding:"required,max=500"`
	ImageIDs []uint `json:"image_ids"`
	VideoURL string `json:"video_url"`
	Location string `json:"location"`
	TopicID  uint   `json:"topic_id"`
}

type UpdateNoteResp struct {
	ID uint `json:"id"`
}

func UpdateNote(c *gin.Context) {
	ctx := c.Request.Context()

	noteIDStr := c.Param("id")
	noteID, err := strconv.ParseUint(noteIDStr, 10, 64)
	if err != nil {
		response.Fail(c, 400, "笔记ID格式错误")
		return
	}

	var req UpdateNoteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "请求参数错误")
		return
	}

	userID := uint(middleware.CurrentUserID(c))
	if userID == 0 {
		response.Fail(c, 401, "未登录")
		return
	}

	note := &model.Note{
		Title:    req.Title,
		Body:     req.Content,
		VideoURL: req.VideoURL,
		Location: req.Location,
		TopicID:  req.TopicID,
		UserID:   userID,
	}
	note.ID = uint(noteID)

	if len(req.ImageIDs) > 0 {
		images, err := imageService.GetByIds(req.ImageIDs)
		if err == nil {
			for _, img := range images {
				note.Images = append(note.Images, model.NoteImage{
					NoteID: note.ID,
					URL:    img.URL,
					Width:  img.Width,
					Height: img.Height,
				})
			}
		}
	}

	if err := noteService.UpdateWithImages(ctx, note); err != nil {
		if err == service.ErrNoteNotFound {
			response.Fail(c, 404, "笔记不存在")
			return
		}
		if err == service.ErrNoPermission {
			response.Fail(c, 403, "无权限")
			return
		}
		response.Fail(c, 500, "更新笔记失败")
		return
	}

	response.Success(c, UpdateNoteResp{ID: note.ID})
}

func DeleteByNoteId(c *gin.Context) {
	ctx := c.Request.Context()

	noteIDStr := c.Param("id")
	noteID, err := strconv.ParseUint(noteIDStr, 10, 64)
	if err != nil {
		response.Fail(c, 400, "笔记ID格式错误")
		return
	}

	userID := uint(middleware.CurrentUserID(c))
	if userID == 0 {
		response.Fail(c, 401, "未登录")
		return
	}

	if err := noteService.DeleteByNoteId(ctx, uint(noteID)); err != nil {
		if err == service.ErrNoteNotFound {
			response.Fail(c, 404, "笔记不存在")
			return
		}
		if err == service.ErrNoPermission {
			response.Fail(c, 403, "无权限")
			return
		}
		response.Fail(c, 500, "删除笔记失败")
		return
	}

	response.Success(c, nil)
}

func SearchNotes(c *gin.Context) {
	ctx := c.Request.Context()

	keyword := c.Query("keyword")
	if keyword == "" {
		response.Fail(c, 400, "搜索关键词不能为空")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	result, err := noteService.SearchNotes(ctx, keyword, page, pageSize)
	if err != nil {
		response.Fail(c, 500, "搜索笔记失败")
		return
	}

	response.Success(c, result)
}

func GetNote(c *gin.Context) {
	ctx := c.Request.Context()

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	result, err := noteService.GetNoteListCached(ctx, page, pageSize)
	if err != nil {
		log.Printf("GetNoteList error: %v", err)
		response.Fail(c, 500, "获取笔记列表失败")
		return
	}

	response.Success(c, result)
}

func GetNoteById(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Fail(c, 400, "笔记ID格式错误")
		return
	}

	note, err := noteService.GetById(ctx, uint(id))
	if err != nil {
		response.Fail(c, 500, "获取笔记详情失败")
		return
	}
	response.Success(c, note)
}

func GetUserNotes(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	userId, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Fail(c, 400, "用户ID格式错误")
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

	notes, total, err := noteService.GetByUserIdWithPagination(ctx, uint(userId), page, pageSize)
	if err != nil {
		response.Fail(c, 500, "获取用户笔记失败")
		return
	}

	totalPage := total / int64(pageSize)
	if total%int64(pageSize) != 0 {
		totalPage++
	}

	response.Success(c, map[string]interface{}{
		"list": notes,
		"pagination": map[string]interface{}{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
			"has_more":  page < int(totalPage),
		},
	})
}

func GetLikedNotes(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	userId, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Fail(c, 400, "用户ID格式错误")
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

	notes, total, err := noteService.GetLikedNotesByUserIdWithPagination(ctx, uint(userId), page, pageSize)
	if err != nil {
		response.Fail(c, 500, "获取赞过列表失败")
		return
	}

	totalPage := total / int64(pageSize)
	if total%int64(pageSize) != 0 {
		totalPage++
	}

	response.Success(c, map[string]interface{}{
		"list": notes,
		"pagination": map[string]interface{}{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
			"has_more":  page < int(totalPage),
		},
	})
}

func LikeNote(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	noteId, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Fail(c, 400, "笔记ID格式错误")
		return
	}

	userId := uint(middleware.CurrentUserID(c))
	if userId == 0 {
		response.Fail(c, 401, "未登录")
		return
	}

	if err := likeService.LikeNote(ctx, userId, uint(noteId)); err != nil {
		response.Fail(c, 500, "点赞失败")
		return
	}
	cache.InvalidateFeed(ctx)
	response.Success(c, "点赞成功")
}

func UnlikeNote(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	noteId, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Fail(c, 400, "笔记ID格式错误")
		return
	}

	userId := uint(middleware.CurrentUserID(c))
	if userId == 0 {
		response.Fail(c, 401, "未登录")
		return
	}

	if err := likeService.UnlikeNote(ctx, userId, uint(noteId)); err != nil {
		response.Fail(c, 500, "取消点赞失败")
		return
	}
	cache.InvalidateFeed(ctx)
	response.Success(c, "取消点赞成功")
}
