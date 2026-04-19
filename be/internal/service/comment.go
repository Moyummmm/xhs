package service

import (
	"errors"
	"server/internal/model"
	"server/internal/repository"
	"time"
)

var (
	ErrCommentNotFound = errors.New("评论不存在")
	ErrInvalidContent  = errors.New("评论内容不能为空")
)

// CommentService 评论服务
type CommentService struct {
	r *repository.CommentRepository
}

func NewCommentService(r *repository.CommentRepository) *CommentService {
	return &CommentService{r: r}
}

// CreateComment 创建评论或回复
func (s *CommentService) CreateComment(userID, noteID uint, content string, parentID uint) error {
	if content == "" {
		return ErrInvalidContent
	}

	comment := &model.Comment{
		Content:  content,
		NoteID:   noteID,
		UserID:   userID,
		ParentID: parentID,
	}

	if parentID > 0 {
		// 回复：更新父评论的 reply_count
		if err := s.r.UpdateReplyCount(parentID, 1); err != nil {
			return err
		}
	} else {
		// 一级评论：更新笔记的 comment_count
		if err := s.r.UpdateNoteCommentCount(noteID, 1); err != nil {
			return err
		}
	}

	return s.r.Create(comment)
}

// DeleteComment 删除评论（仅允许删除自己的评论）
func (s *CommentService) DeleteComment(userID, commentID uint) error {
	comment, err := s.r.GetByID(commentID)
	if err != nil {
		return ErrCommentNotFound
	}

	if comment.UserID != userID {
		return ErrNoPermission
	}

	if err := s.r.Delete(commentID); err != nil {
		return err
	}

	if comment.ParentID > 0 {
		// 回复：减少父评论的 reply_count
		_ = s.r.UpdateReplyCount(comment.ParentID, -1)
	} else {
		// 一级评论：减少笔记的 comment_count
		_ = s.r.UpdateNoteCommentCount(comment.NoteID, -1)
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

// GetComments 获取笔记的评论列表
func (s *CommentService) GetComments(noteID uint, page, pageSize int, sort string, currentUserID uint) (*CommentListResp, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 50 {
		pageSize = 50
	}

	comments, total, err := s.r.GetByNoteId(noteID, page, pageSize, sort)
	if err != nil {
		return nil, err
	}

	// 收集所有一级评论 ID
	parentIDs := make([]uint, 0, len(comments))
	for _, c := range comments {
		parentIDs = append(parentIDs, c.ID)
	}

	// 批量获取回复
	repliesMap := make(map[uint][]model.Comment)
	if len(parentIDs) > 0 {
		allReplies, err := s.r.GetRepliesByParentIDs(parentIDs)
		if err == nil {
			for _, r := range allReplies {
				repliesMap[r.ParentID] = append(repliesMap[r.ParentID], r)
			}
		}
	}

	// 批量查询点赞状态
	likedSet := make(map[uint]bool)
	if currentUserID > 0 && len(parentIDs) > 0 {
		for _, c := range comments {
			isLiked, _ := s.r.IsLikedByUser(currentUserID, c.ID)
			if isLiked {
				likedSet[c.ID] = true
			}
		}
		// 也查询回复的点赞状态（一级评论展示的前3条回复）
		for _, replies := range repliesMap {
			for _, r := range replies {
				isLiked, _ := s.r.IsLikedByUser(currentUserID, r.ID)
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

		// 每条一级评论最多展示 3 条回复
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
