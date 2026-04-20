package repository

import (
	"context"
	"server/internal/model"

	"gorm.io/gorm"
)

type LikeRepository struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) *LikeRepository {
	return &LikeRepository{db: db}
}

func (r *LikeRepository) LikeNote(ctx context.Context, userId, noteId uint) error {
	var like model.Like4Note

	err := r.db.WithContext(ctx).Unscoped().Where("note_id = ? AND user_id = ?", noteId, userId).First(&like).Error
	if err == nil {
		if like.DeletedAt.Valid {
			err = r.db.WithContext(ctx).Unscoped().Model(&like).Update("deleted_at", nil).Error
			if err != nil {
				return err
			}
		}
		return nil
	}

	if err == gorm.ErrRecordNotFound {
		like = model.Like4Note{
			NoteID: noteId,
			UserID: userId,
		}
		if err := r.db.WithContext(ctx).Save(&like).Error; err != nil {
			return err
		}
	} else {
		return err
	}

	return r.db.WithContext(ctx).Model(&model.Note{}).Where("id = ?", noteId).
		UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error
}

func (r *LikeRepository) UnlikeNote(ctx context.Context, userId, noteId uint) error {
	if err := r.db.WithContext(ctx).Where("user_id = ? AND note_id = ?", userId, noteId).
		Delete(&model.Like4Note{}).Error; err != nil {
		return err
	}
	return r.db.WithContext(ctx).Model(&model.Note{}).Where("id = ?", noteId).
		UpdateColumn("like_count", gorm.Expr("like_count - ?", 1)).Error
}

func (r *LikeRepository) GetLikeCountByNoteId(ctx context.Context, noteId uint) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Like4Note{}).Where("note_id = ?", noteId).Count(&count).Error
	return int(count), err
}

func (r *LikeRepository) IsLiked(ctx context.Context, userId, noteId uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Like4Note{}).Where("user_id = ? AND note_id = ?", userId, noteId).Count(&count).Error
	return count > 0, err
}
