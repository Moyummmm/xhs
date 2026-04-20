package repository

import (
	"context"
	"server/internal/model"

	"gorm.io/gorm"
)

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Create(ctx context.Context, comment *model.Comment) error {
	if comment.ParentID == nil {
		// 顶级评论：手动插入 NULL 而非 0，避免外键约束冲突
		return r.db.WithContext(ctx).Exec("INSERT INTO comments (content, note_id, user_id, parent_id, like_count, reply_count, created_at, updated_at) VALUES (?, ?, ?, NULL, ?, ?, NOW(), NOW())",
			comment.Content, comment.NoteID, comment.UserID, comment.LikeCount, comment.ReplyCount).Error
	}
	return r.db.WithContext(ctx).Create(comment).Error
}

func (r *CommentRepository) Delete(ctx context.Context, commentID uint) error {
	return r.db.WithContext(ctx).Delete(&model.Comment{}, commentID).Error
}

func (r *CommentRepository) GetByID(ctx context.Context, commentID uint) (*model.Comment, error) {
	var comment model.Comment
	err := r.db.WithContext(ctx).First(&comment, commentID).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *CommentRepository) GetByNoteId(ctx context.Context, noteID uint, page, pageSize int, sort string) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.WithContext(ctx).Model(&model.Comment{}).Where("note_id = ? AND parent_id IS NULL", noteID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db := query.Preload("User").Offset(offset).Limit(pageSize)
	if sort == "latest" {
		db = db.Order("created_at DESC")
	} else {
		db = db.Order("like_count DESC, created_at DESC")
	}

	err := db.Find(&comments).Error
	return comments, total, err
}

func (r *CommentRepository) GetRepliesByParentIDs(ctx context.Context, parentIDs []uint) ([]model.Comment, error) {
	if len(parentIDs) == 0 {
		return nil, nil
	}

	var comments []model.Comment
	err := r.db.WithContext(ctx).Where("parent_id IN ?", parentIDs).
		Preload("User").
		Order("parent_id ASC, created_at DESC").
		Find(&comments).Error
	return comments, err
}

func (r *CommentRepository) CountReplies(ctx context.Context, parentID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Comment{}).Where("parent_id = ?", parentID).Count(&count).Error
	return count, err
}

func (r *CommentRepository) UpdateReplyCount(ctx context.Context, commentID uint, delta int) error {
	return r.db.WithContext(ctx).Model(&model.Comment{}).Where("id = ?", commentID).
		UpdateColumn("reply_count", gorm.Expr("reply_count + ?", delta)).Error
}

func (r *CommentRepository) UpdateNoteCommentCount(ctx context.Context, noteID uint, delta int) error {
	return r.db.WithContext(ctx).Model(&model.Note{}).Where("id = ?", noteID).
		UpdateColumn("comment_count", gorm.Expr("comment_count + ?", delta)).Error
}

func (r *CommentRepository) IsLikedByUser(ctx context.Context, userID, commentID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Like4Comment{}).Where("user_id = ? AND comment_id = ?", userID, commentID).Count(&count).Error
	return count > 0, err
}
