package repository

import (
	"context"
	"server/internal/model"

	"gorm.io/gorm"
)

type CollectRepository struct {
	db *gorm.DB
}

func NewCollectRepository(db *gorm.DB) *CollectRepository {
	return &CollectRepository{db: db}
}

func (r *CollectRepository) CollectById(ctx context.Context, userId, noteId uint) error {
	var collection model.Collect

	err := r.db.WithContext(ctx).Unscoped().Where("note_id = ? AND user_id = ?", noteId, userId).First(&collection).Error
	if err == nil {
		if collection.DeletedAt.Valid {
			err = r.db.WithContext(ctx).Unscoped().Model(&collection).Update("deleted_at", nil).Error
			if err != nil {
				return err
			}
		}
		return nil
	}

	if err == gorm.ErrRecordNotFound {
		collection = model.Collect{
			NoteID: noteId,
			UserID: userId,
			Active: true,
		}
		if err := r.db.WithContext(ctx).Save(&collection).Error; err != nil {
			return err
		}
	} else {
		return err
	}

	return nil
}

func (r *CollectRepository) DisCollectById(ctx context.Context, userId, noteId uint) error {
	return r.db.WithContext(ctx).Where("note_id = ? AND user_id = ?", noteId, userId).
		Delete(&model.Collect{}).Error
}

func (r *CollectRepository) GetCollecCountsByUserId(ctx context.Context, userId uint) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Collect{}).Where("user_id = ?", userId).Count(&count).Error
	return int(count), err
}

func (r *CollectRepository) GetCollectList(ctx context.Context, userId uint) ([]model.Note, error) {
	var notes []model.Note
	err := r.db.WithContext(ctx).Joins("JOIN collects ON collects.note_id = notes.id").
		Where("collects.user_id = ? AND collects.active = ?", userId, true).
		Find(&notes).Error
	return notes, err
}

func (r *CollectRepository) GetCollectListWithPagination(ctx context.Context, userId uint, page, pageSize int) ([]model.Note, int64, error) {
	var notes []model.Note
	var total int64

	offset := (page - 1) * pageSize

	query := r.db.WithContext(ctx).Model(&model.Note{}).
		Joins("JOIN collects ON collects.note_id = notes.id").
		Where("collects.user_id = ? AND collects.active = ?", userId, true)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("User").Preload("Images").
		Order("collects.created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&notes).Error

	return notes, total, err
}
