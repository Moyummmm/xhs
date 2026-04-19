package repository

import (
	"server/internal/model"

	"gorm.io/gorm"
)

type CollectRepository struct {
	db *gorm.DB
}

func NewCollectRepository(db *gorm.DB) *CollectRepository {
	return &CollectRepository{db: db}
}

// CollectById 收藏笔记
// userId: 用户ID
// noteId: 笔记ID
func (r *CollectRepository) CollectById(userId, noteId uint) error {
	collection := model.Collect{
		NoteID: noteId,
		UserID: userId,
	}
	return r.db.Save(&collection).Error
}

// DisCollectById 取消收藏
// userId: 用户ID
// noteId: 笔记ID
func (r *CollectRepository) DisCollectById(userId, noteId uint) error {
	collection := model.Collect{
		NoteID: noteId,
		UserID: userId,
		Active: false,
	}
	return r.db.Save(&collection).Error
}

// GetCollecCountsByUserId 获取用户收藏数量
// userId: 用户ID
func (r *CollectRepository) GetCollecCountsByUserId(userId uint) (int, error) {
	var count int64
	err := r.db.Model(&model.Collect{}).Where("user_id = ?", userId).Count(&count).Error
	return int(count), err
}

// GetCollectList 获取用户收藏列表
// userId: 用户ID
func (r *CollectRepository) GetCollectList(userId uint) ([]model.Note, error) {
	var notes []model.Note
	err := r.db.Joins("JOIN collects ON collects.note_id = notes.id").
		Where("collects.user_id = ? AND collects.active = ?", userId, true).
		Find(&notes).Error
	return notes, err
}
