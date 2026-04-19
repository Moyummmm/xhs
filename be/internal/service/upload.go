package service

import (
	"context"
	"fmt"
	"image"
	"mime/multipart"

	"server/config"
)

type UploadService struct{}

func NewUploadService() *UploadService {
	return &UploadService{}
}

type ImageUploadResult struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type VideoUploadResult struct {
	URL      string `json:"url"`
	Duration int    `json:"duration"`
	CoverURL string `json:"cover_url"`
}

func (s *UploadService) UploadImage(ctx context.Context, file *multipart.FileHeader) (*ImageUploadResult, error) {
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("open file failed: %v", err)
	}
	defer src.Close()

	url, err := config.UploadFileReader(ctx, file.Filename, src, file.Size, file.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}

	width, height := getImageDimensions(src)

	return &ImageUploadResult{
		URL:    url,
		Width:  width,
		Height: height,
	}, nil
}

func (s *UploadService) UploadVideo(ctx context.Context, file *multipart.FileHeader) (*VideoUploadResult, error) {
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("open file failed: %v", err)
	}
	defer src.Close()

	url, err := config.UploadFileReader(ctx, file.Filename, src, file.Size, file.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}

	return &VideoUploadResult{
		URL:      url,
		Duration: 0,
		CoverURL: "",
	}, nil
}

func getImageDimensions(src multipart.File) (int, int) {
	cfg, _, err := image.DecodeConfig(src)
	if err != nil {
		return 0, 0
	}
	return cfg.Width, cfg.Height
}
