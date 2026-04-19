package repository

import (
	"server/internal/model"

	"gorm.io/gorm"
)

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

// Create 创建评论
func (r *CommentRepository) Create(comment *model.Comment) error {
	return r.db.Create(comment).Error
}

// Delete 软删除评论
func (r *CommentRepository) Delete(commentID uint) error {
	return r.db.Delete(&model.Comment{}, commentID).Error
}

// GetByID 根据ID获取评论
func (r *CommentRepository) GetByID(commentID uint) (*model.Comment, error) {
	var comment model.Comment
	err := r.db.First(&comment, commentID).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// GetByNoteId 分页获取一级评论列表（parent_id = 0）
func (r *CommentRepository) GetByNoteId(noteID uint, page, pageSize int, sort string) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&model.Comment{}).Where("note_id = ? AND parent_id = 0", noteID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db := query.Preload("User").Offset(offset).Limit(pageSize)
	if sort == "latest" {
		db = db.Order("created_at DESC")
	} else {
		// hot: 按点赞数降序，再按时间降序
		db = db.Order("like_count DESC, created_at DESC")
	}

	err := db.Find(&comments).Error
	return comments, total, err
}

// GetRepliesByParentIDs 批量获取指定一级评论的回复
// 返回所有回复，由调用方按 parent_id 分组并截取前 limit 条
func (r *CommentRepository) GetRepliesByParentIDs(parentIDs []uint) ([]model.Comment, error) {
	if len(parentIDs) == 0 {
		return nil, nil
	}

	var comments []model.Comment
	err := r.db.Where("parent_id IN ?", parentIDs).
		Preload("User").
		Order("parent_id ASC, created_at DESC").
		Find(&comments).Error
	return comments, err
}

// CountReplies 统计某一级评论的回复数
func (r *CommentRepository) CountReplies(parentID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Comment{}).Where("parent_id = ?", parentID).Count(&count).Error
	return count, err
}

// UpdateReplyCount 原子更新评论的 reply_count
func (r *CommentRepository) UpdateReplyCount(commentID uint, delta int) error {
	return r.db.Model(&model.Comment{}).Where("id = ?", commentID).
		UpdateColumn("reply_count", gorm.Expr("reply_count + ?", delta)).Error
}

// UpdateNoteCommentCount 原子更新笔记的 comment_count
func (r *CommentRepository) UpdateNoteCommentCount(noteID uint, delta int) error {
	return r.db.Model(&model.Note{}).Where("id = ?", noteID).
		UpdateColumn("comment_count", gorm.Expr("comment_count + ?", delta)).Error
}

// IsLikedByUser 检查用户是否点赞了某条评论
func (r *CommentRepository) IsLikedByUser(userID, commentID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.Like4Comment{}).Where("user_id = ? AND comment_id = ?", userID, commentID).Count(&count).Error
	return count > 0, err
}
