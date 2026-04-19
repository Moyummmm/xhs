package handler

import (
	"log"
	"server/config"
	"server/internal/middleware"
	"server/internal/model"
	"server/internal/repository"
	"server/internal/service"
	"server/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// noteService 笔记服务实例
var noteService *service.NoteService
var likeService *service.LikeService
// imageService 已由 upload.go 的 init() 初始化，共用同一个实例

// init 初始化函数，创建笔记服务和仓储实例
func init() {
	noteRepository := repository.NewNoteRepository(config.DB())
	noteService = service.NewNoteService(noteRepository)

	likeRepo := repository.NewLikeRepository(config.DB())
	likeService = service.NewLikeService(likeRepo)
}

// CreateNoteReq 创建笔记请求结构体
type CreateNoteReq struct {
	Title    string `json:"title" binding:"required,max=100"`
	Content  string `json:"content" binding:"required,max=500"`
	ImageIDs []uint `json:"image_ids"`
	VideoURL string `json:"video_url"`
	Location string `json:"location"`
	TopicID  uint   `json:"topic_id"`
}

// CreateNoteResp 创建笔记响应结构体
type CreateNoteResp struct {
	ID uint `json:"id"` // 笔记ID
}

// Create 创建笔记接口
// @Summary 创建笔记
// @Description 创建一条新笔记，支持添加图片
// @Tags 笔记
// @Accept json
// @Produce json
// @Param request body CreateNoteReq true "创建笔记请求参数"
// @Success 200 {object} response.Response{data=CreateNoteResp} "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未登录"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/notes [post]
func CreateNote(c *gin.Context) {
	// 绑定请求参数
	var req CreateNoteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "请求参数错误")
		return
	}

	// 获取当前登录用户ID
	userID := uint(middleware.CurrentUserID(c))
	if userID == 0 {
		response.Fail(c, 401, "未登录")
		return
	}

	// 构建笔记模型
	note := &model.Note{
		Title:    req.Title,
		Body:     req.Content,
		VideoURL: req.VideoURL,
		Location: req.Location,
		TopicID:  req.TopicID,
		UserID:   userID,
	}

	// 处理图片列表
	if len(req.ImageIDs) > 0 {
		images, err := imageService.GetByIds(req.ImageIDs)
		if err == nil {
			for _, img := range images {
				note.Images = append(note.Images, model.NoteImage{
					NoteID: note.ID,
					URL:    img.URL,
				})
			}
		}
	}

	// 调用服务创建笔记
	if err := noteService.CreateWithImages(note); err != nil {
		response.Fail(c, 500, "创建笔记失败")
		return
	}

	// 返回创建成功的响应
	response.Success(c, CreateNoteResp{ID: note.ID})
}

// UpdateNoteReq 更新笔记请求结构体
type UpdateNoteReq struct {
	Title    string   `json:"title" binding:"required,max=100"`
	Content  string   `json:"content" binding:"required,max=500"`
	ImageIDs []uint   `json:"image_ids"`
	VideoURL string   `json:"video_url"`
	Location string   `json:"location"`
	TopicID  uint     `json:"topic_id"`
}

// UpdateNoteResp 更新笔记响应结构体
type UpdateNoteResp struct {
	ID uint `json:"id"`
}

func UpdateNote(c *gin.Context) {
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
				})
			}
		}
	}

	if err := noteService.UpdateWithImages(note); err != nil {
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

// DeleteByNoteId 依据NoteId删除Note
func DeleteByNoteId(c *gin.Context) {
	// 获取笔记ID参数
	noteIDStr := c.Param("id")
	noteID, err := strconv.ParseUint(noteIDStr, 10, 64)
	if err != nil {
		response.Fail(c, 400, "笔记ID格式错误")
		return
	}

	// 获取当前登录用户ID
	userID := uint(middleware.CurrentUserID(c))
	if userID == 0 {
		response.Fail(c, 401, "未登录")
		return
	}

	// 调用服务删除笔记
	if err := noteService.DeleteByNoteId(uint(noteID)); err != nil {
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

	// 返回成功响应
	response.Success(c, nil)
}

// SearchNotes 搜索笔记
// GET /notes/search
func SearchNotes(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		response.Fail(c, 400, "搜索关键词不能为空")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	result, err := noteService.SearchNotes(keyword, page, pageSize)
	if err != nil {
		response.Fail(c, 500, "搜索笔记失败")
		return
	}

	response.Success(c, result)
}

// GetNote 获取笔记
func GetNote(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	result, err := noteService.GetNoteList(page, pageSize)
	if err != nil {
		log.Printf("GetNoteList error: %v", err)
		response.Fail(c, 500, "获取笔记列表失败")
		return
	}

	response.Success(c, result)
}

// GetNoteById 获取笔记详情
// GET /notes/:id
func GetNoteById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Fail(c, 400, "笔记ID格式错误")
		return
	}

	note, err := noteService.GetById(uint(id))
	if err != nil {
		response.Fail(c, 500, "获取笔记详情失败")
		return
	}
	response.Success(c, note)
}

// GetUserNotes 获取用户笔记列表
// GET /users/:id/notes
func GetUserNotes(c *gin.Context) {
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

	notes, total, err := noteService.GetByUserIdWithPagination(uint(userId), page, pageSize)
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

// LikeNote 点赞笔记
// POST /notes/:id/like
func LikeNote(c *gin.Context) {
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

	if err := likeService.LikeNote(userId, uint(noteId)); err != nil {
		response.Fail(c, 500, "点赞失败")
		return
	}
	response.Success(c, "点赞成功")
}

// UnlikeNote 取消点赞笔记
// DELETE /notes/:id/like
func UnlikeNote(c *gin.Context) {
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

	if err := likeService.UnlikeNote(userId, uint(noteId)); err != nil {
		response.Fail(c, 500, "取消点赞失败")
		return
	}
	response.Success(c, "取消点赞成功")
}
