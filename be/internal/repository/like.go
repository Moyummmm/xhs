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

// LikeNote 使用 CTE + RETURNING 确保只有在真正新增或恢复已删除点赞时才 +1
// 避免了重复点击已存在的有效点赞导致的无故 +1
func (r *LikeRepository) LikeNote(ctx context.Context, userId, noteId uint) error {
	sql := `
		WITH upserted AS (
			INSERT INTO like4_notes (note_id, user_id, deleted_at, updated_at)
			VALUES ($1, $2, NULL, NOW())
			ON CONFLICT (note_id, user_id)
			DO UPDATE SET
				deleted_at = NULL,
				updated_at = NOW()
			WHERE like4_notes.deleted_at IS NOT NULL
			RETURNING 1 as updated
		)
		UPDATE notes
		SET like_count = like_count + 1
		WHERE id = $3
		  AND EXISTS (SELECT 1 FROM upserted)
	`
	result := r.db.WithContext(ctx).Exec(sql, noteId, userId, noteId)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// UnlikeNote 使用 CTE + 条件 UPDATE 确保只有在真正删除了有效点赞时才 -1
// 避免并发双击导致重复扣减 like_count
func (r *LikeRepository) UnlikeNote(ctx context.Context, userId, noteId uint) error {
	sql := `
		WITH deleted AS (
			UPDATE like4_notes
			SET deleted_at = NOW()
			WHERE note_id = $1
			  AND user_id = $2
			  AND deleted_at IS NULL
			RETURNING 1
		)
		UPDATE notes
		SET like_count = GREATEST(like_count - 1, 0)
		WHERE id = $3
		  AND EXISTS (SELECT 1 FROM deleted)
	`
	result := r.db.WithContext(ctx).Exec(sql, noteId, userId, noteId)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *LikeRepository) GetLikeCountByNoteId(ctx context.Context, noteId uint) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Like4Note{}).Where("note_id = ?", noteId).Count(&count).Error
	return int(count), err
}

func (r *LikeRepository) IsLiked(ctx context.Context, userId, noteId uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Like4Note{}).Where("user_id = ? AND note_id = ?", noteId, userId).Count(&count).Error
	return count > 0, err
}