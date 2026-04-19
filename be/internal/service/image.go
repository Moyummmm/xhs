package service

import (
	"context"

	"server/internal/model"
	"server/internal/repository"
)

type ImageService struct {
	r *repository.ImageRepository
}

func NewImageService(r *repository.ImageRepository) *ImageService {
	return &ImageService{r: r}
}

func (s *ImageService) Create(ctx context.Context, url string, width, height int, userId uint) (*model.Image, error) {
	img := &model.Image{
		URL:    url,
		Width:  width,
		Height: height,
		UserID: userId,
	}
	if err := s.r.Create(img); err != nil {
		return nil, err
	}
	return img, nil
}

func (s *ImageService) GetById(id uint) (*model.Image, error) {
	return s.r.GetById(id)
}

func (s *ImageService) GetByIds(ids []uint) ([]model.Image, error) {
	return s.r.GetByIds(ids)
}

func (s *ImageService) DeleteById(id uint) error {
	return s.r.DeleteById(id)
}
