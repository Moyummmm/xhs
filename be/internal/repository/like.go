package repository

import (
	"server/internal/model"

	"gorm.io/gorm"
)

type LikeRepository struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) *LikeRepository {
	return &LikeRepository{db: db}
}

// LikeNote 点赞笔记
func (r *LikeRepository) LikeNote(userId, noteId uint) error {
	like := model.Like4Note{
		NoteID: noteId,
		UserID: userId,
	}
	if err := r.db.Save(&like).Error; err != nil {
		return err
	}
	return r.db.Model(&model.Note{}).Where("id = ?", noteId).
		UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error
}

// UnlikeNote 取消点赞笔记
func (r *LikeRepository) UnlikeNote(userId, noteId uint) error {
	if err := r.db.Where("user_id = ? AND note_id = ?", userId, noteId).
		Delete(&model.Like4Note{}).Error; err != nil {
		return err
	}
	return r.db.Model(&model.Note{}).Where("id = ?", noteId).
		UpdateColumn("like_count", gorm.Expr("like_count - ?", 1)).Error
}

// GetLikeCountByNoteId 获取笔记点赞数量
func (r *LikeRepository) GetLikeCountByNoteId(noteId uint) (int, error) {
	var count int64
	err := r.db.Model(&model.Like4Note{}).Where("note_id = ?", noteId).Count(&count).Error
	return int(count), err
}

// IsLiked 检查用户是否点赞了笔记
func (r *LikeRepository) IsLiked(userId, noteId uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.Like4Note{}).Where("user_id = ? AND note_id = ?", userId, noteId).Count(&count).Error
	return count > 0, err
}
