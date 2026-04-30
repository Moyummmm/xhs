package service

import (
	"context"
	"errors"
	"server/internal/model"
	"server/internal/repository"
	"time"

	"gorm.io/gorm"
)

var (
	ErrCommentNotFound = errors.New("评论不存在")
	ErrInvalidContent  = errors.New("评论内容不能为空")
)

type CommentService struct {
	r   *repository.CommentRepository
	db  *gorm.DB
}

func NewCommentService(r *repository.CommentRepository, db *gorm.DB) *CommentService {
	return &CommentService{r: r, db: db}
}

// CreateComment 创建评论，使用事务确保评论创建和计数更新原子性
func (s *CommentService) CreateComment(ctx context.Context, userID, noteID uint, content string, parentID *uint) error {
	if content == "" {
		return ErrInvalidContent
	}

	if s.db == nil {
		return errors.New("database not initialized")
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		comment := &model.Comment{
			Content:  content,
			NoteID:   noteID,
			UserID:   userID,
			ParentID: parentID,
		}

		// 先创建评论
		if parentID != nil && *parentID > 0 {
			// 创建回复：需要在事务内执行创建，并更新父评论的回复数
			if err := tx.Create(comment).Error; err != nil {
				return err
			}
			// 更新父评论的回复数
			if err := tx.Model(&model.Comment{}).Where("id = ?", *parentID).
				UpdateColumn("reply_count", gorm.Expr("reply_count + ?", 1)).Error; err != nil {
				return err
			}
		} else {
			// 创建顶级评论：使用原生 SQL 插入 NULL 而非 0
			if err := tx.Exec("INSERT INTO comments (content, note_id, user_id, parent_id, like_count, reply_count, created_at, updated_at) VALUES (?, ?, ?, NULL, ?, ?, NOW(), NOW())",
				content, noteID, userID, 0, 0).Error; err != nil {
				return err
			}
			// 更新笔记的评论数
			if err := tx.Model(&model.Note{}).Where("id = ?", noteID).
				UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1)).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// DeleteComment 删除评论，使用事务确保评论删除和计数更新原子性
func (s *CommentService) DeleteComment(ctx context.Context, userID, commentID uint) error {
	if s.db == nil {
		return errors.New("database not initialized")
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先查询评论（使用 FOR UPDATE 防止并发问题）
		var comment model.Comment
		if err := tx.Set("gorm:query_option", "FOR UPDATE").
			First(&comment, commentID).Error; err != nil {
			return ErrCommentNotFound
		}

		if comment.UserID != userID {
			return ErrNoPermission
		}

		// 删除评论
		if err := tx.Delete(&model.Comment{}, commentID).Error; err != nil {
			return err
		}

		// 更新计数
		if comment.ParentID != nil && *comment.ParentID > 0 {
			// 回复评论：减少父评论的回复数
			if err := tx.Model(&model.Comment{}).Where("id = ?", *comment.ParentID).
				UpdateColumn("reply_count", gorm.Expr("GREATEST(reply_count - 1, 0)")).Error; err != nil {
				return err
			}
		} else {
			// 顶级评论：减少笔记的评论数
			if err := tx.Model(&model.Note{}).Where("id = ?", comment.NoteID).
				UpdateColumn("comment_count", gorm.Expr("GREATEST(comment_count - 1, 0)")).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// CommentItem 评论列表项（一级评论）
type CommentItem struct {
	ID         uint        `json:"id"`
	User       model.User  `json:"user"`
	Content    string      `json:"content"`
	LikeCount  uint        `json:"like_count"`
	ReplyCount uint        `json:"reply_count"`
	IsLiked    bool        `json:"is_liked"`
	Replies    []ReplyItem `json:"replies"`
	CreatedAt  time.Time   `json:"created_at"`
}

// ReplyItem 回复项
type ReplyItem struct {
	ID        uint       `json:"id"`
	User      model.User `json:"user"`
	Content   string     `json:"content"`
	LikeCount uint       `json:"like_count"`
	CreatedAt time.Time  `json:"created_at"`
}

// CommentListResp 评论列表响应
type CommentListResp struct {
	List     []CommentItem `json:"list"`
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
}

func (s *CommentService) GetComments(ctx context.Context, noteID uint, page, pageSize int, sort string, currentUserID uint) (*CommentListResp, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 50 {
		pageSize = 50
	}

	comments, total, err := s.r.GetByNoteId(ctx, noteID, page, pageSize, sort)
	if err != nil {
		return nil, err
	}

	parentIDs := make([]uint, 0, len(comments))
	for _, c := range comments {
		parentIDs = append(parentIDs, c.ID)
	}

	repliesMap := make(map[uint][]model.Comment)
	if len(parentIDs) > 0 {
		allReplies, err := s.r.GetRepliesByParentIDs(ctx, parentIDs)
		if err == nil {
			for _, r := range allReplies {
				if r.ParentID != nil {
					repliesMap[*r.ParentID] = append(repliesMap[*r.ParentID], r)
				}
			}
		}
	}

	likedSet := make(map[uint]bool)
	if currentUserID > 0 && len(parentIDs) > 0 {
		allCommentIDs := make([]uint, 0, len(comments))
		for _, c := range comments {
			allCommentIDs = append(allCommentIDs, c.ID)
		}
		for _, replies := range repliesMap {
			for _, r := range replies {
				allCommentIDs = append(allCommentIDs, r.ID)
			}
		}
		batchLiked, _ := s.r.BatchIsLikedByUser(ctx, currentUserID, allCommentIDs)
		for id, v := range batchLiked {
			likedSet[id] = v
		}
	}

	list := make([]CommentItem, 0, len(comments))
	for _, c := range comments {
		item := CommentItem{
			ID:         c.ID,
			User:       c.User,
			Content:    c.Content,
			LikeCount:  c.LikeCount,
			ReplyCount: c.ReplyCount,
			IsLiked:    likedSet[c.ID],
			CreatedAt:  c.CreatedAt,
			Replies:    make([]ReplyItem, 0),
		}

		replies := repliesMap[c.ID]
		limit := 3
		if len(replies) < limit {
			limit = len(replies)
		}
		for i := 0; i < limit; i++ {
			r := replies[i]
			item.Replies = append(item.Replies, ReplyItem{
				ID:        r.ID,
				User:      r.User,
				Content:   r.Content,
				LikeCount: r.LikeCount,
				CreatedAt: r.CreatedAt,
			})
		}

		list = append(list, item)
	}

	return &CommentListResp{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}