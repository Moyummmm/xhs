package repository

import (
	"server/internal/model"

	"gorm.io/gorm"
)

type ImageRepository struct {
	db *gorm.DB
}

func NewImageRepository(db *gorm.DB) *ImageRepository {
	return &ImageRepository{db: db}
}

func (r *ImageRepository) Create(image *model.Image) error {
	return r.db.Create(image).Error
}

func (r *ImageRepository) GetById(id uint) (*model.Image, error) {
	var image model.Image
	if err := r.db.Where("id = ?", id).First(&image).Error; err != nil {
		return nil, err
	}
	return &image, nil
}

func (r *ImageRepository) GetByIds(ids []uint) ([]model.Image, error) {
	var images []model.Image
	if err := r.db.Where("id IN ?", ids).Find(&images).Error; err != nil {
		return nil, err
	}
	return images, nil
}

func (r *ImageRepository) DeleteById(id uint) error {
	return r.db.Where("id = ?", id).Delete(&model.Image{}).Error
}
