package service

import (
	"context"
	"errors"
	"server/internal/model"
	"server/internal/repository"
	"time"
)

var (
	ErrCommentNotFound = errors.New("评论不存在")
	ErrInvalidContent  = errors.New("评论内容不能为空")
)

type CommentService struct {
	r *repository.CommentRepository
}

func NewCommentService(r *repository.CommentRepository) *CommentService {
	return &CommentService{r: r}
}

func (s *CommentService) CreateComment(ctx context.Context, userID, noteID uint, content string, parentID *uint) error {
	if content == "" {
		return ErrInvalidContent
	}

	comment := &model.Comment{
		Content:  content,
		NoteID:   noteID,
		UserID:   userID,
		ParentID: parentID,
	}

	if parentID != nil && *parentID > 0 {
		if err := s.r.UpdateReplyCount(ctx, *parentID, 1); err != nil {
			return err
		}
	} else {
		if err := s.r.UpdateNoteCommentCount(ctx, noteID, 1); err != nil {
			return err
		}
	}

	return s.r.Create(ctx, comment)
}

func (s *CommentService) DeleteComment(ctx context.Context, userID, commentID uint) error {
	comment, err := s.r.GetByID(ctx, commentID)
	if err != nil {
		return ErrCommentNotFound
	}

	if comment.UserID != userID {
		return ErrNoPermission
	}

	if err := s.r.Delete(ctx, commentID); err != nil {
		return err
	}

	if comment.ParentID != nil && *comment.ParentID > 0 {
		_ = s.r.UpdateReplyCount(ctx, *comment.ParentID, -1)
	} else {
		_ = s.r.UpdateNoteCommentCount(ctx, comment.NoteID, -1)
	}

	return nil
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
		for _, c := range comments {
			isLiked, _ := s.r.IsLikedByUser(ctx, currentUserID, c.ID)
			if isLiked {
				likedSet[c.ID] = true
			}
		}
		for _, replies := range repliesMap {
			for _, r := range replies {
				isLiked, _ := s.r.IsLikedByUser(ctx, currentUserID, r.ID)
				if isLiked {
					likedSet[r.ID] = true
				}
			}
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
