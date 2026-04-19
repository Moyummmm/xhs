package handler

import (
	"server/config"
	"server/internal/middleware"
	"server/internal/repository"
	"server/internal/service"
	"server/pkg/response"

	"github.com/gin-gonic/gin"
)

var uploadService *service.UploadService
var imageService *service.ImageService

func init() {
	uploadService = service.NewUploadService()
	imageRepo := repository.NewImageRepository(config.DB())
	imageService = service.NewImageService(imageRepo)
}

// UploadImage 上传图片
// POST /upload/image
func UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, 400, "请选择要上传的文件")
		return
	}

	result, err := uploadService.UploadImage(c.Request.Context(), file)
	if err != nil {
		response.Fail(c, 500, "上传图片失败")
		return
	}

	userId := uint(middleware.CurrentUserID(c))

	img, err := imageService.Create(c.Request.Context(), result.URL, result.Width, result.Height, userId)
	if err != nil {
		response.Fail(c, 500, "保存图片记录失败")
		return
	}

	response.Success(c, map[string]interface{}{
		"id":      img.ID,
		"url":     img.URL,
		"width":   img.Width,
		"height":  img.Height,
	})
}

// UploadVideo 上传视频
// POST /upload/video
func UploadVideo(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, 400, "请选择要上传的文件")
		return
	}

	result, err := uploadService.UploadVideo(c.Request.Context(), file)
	if err != nil {
		response.Fail(c, 500, "上传视频失败")
		return
	}

	userId := uint(middleware.CurrentUserID(c))

	img, err := imageService.Create(c.Request.Context(), result.URL, 0, 0, userId)
	if err != nil {
		response.Fail(c, 500, "保存视频记录失败")
		return
	}

	response.Success(c, map[string]interface{}{
		"id":        img.ID,
		"url":       img.URL,
		"duration":  result.Duration,
		"cover_url": result.CoverURL,
	})
}
