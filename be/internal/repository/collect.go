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

// CollectById 使用 CTE + RETURNING 确保只有在真正新增或恢复已删除收藏时才 +1
// 避免了重复点击已存在的有效收藏导致的无故 +1
func (r *CollectRepository) CollectById(ctx context.Context, userId, noteId uint) error {
	sql := `
		WITH upserted AS (
			INSERT INTO collects (note_id, user_id, active, deleted_at, updated_at)
			VALUES ($1, $2, true, NULL, NOW())
			ON CONFLICT (note_id, user_id)
			DO UPDATE SET
				deleted_at = NULL,
				active = true,
				updated_at = NOW()
			WHERE collects.deleted_at IS NOT NULL
			RETURNING 1 as updated
		)
		UPDATE notes
		SET collect_count = collect_count + 1
		WHERE id = $3
		  AND EXISTS (SELECT 1 FROM upserted)
	`
	result := r.db.WithContext(ctx).Exec(sql, noteId, userId, noteId)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// DisCollectById 使用 CTE + 条件 UPDATE 确保只有在真正删除了有效收藏时才 -1
// 避免并发双击导致重复扣减 collect_count
func (r *CollectRepository) DisCollectById(ctx context.Context, userId, noteId uint) error {
	sql := `
		WITH deleted AS (
			UPDATE collects
			SET deleted_at = NOW()
			WHERE note_id = $1
			  AND user_id = $2
			  AND deleted_at IS NULL
			RETURNING 1
		)
		UPDATE notes
		SET collect_count = GREATEST(collect_count - 1, 0)
		WHERE id = $3
		  AND EXISTS (SELECT 1 FROM deleted)
	`
	result := r.db.WithContext(ctx).Exec(sql, noteId, userId, noteId)
	if result.Error != nil {
		return result.Error
	}
	return nil
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