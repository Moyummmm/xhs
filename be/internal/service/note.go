package service

import (
	"context"
	"errors"
	"server/internal/model"
	"server/internal/repository"
)

var (
	ErrNoteNotFound = errors.New("笔记不存在")
	ErrNoPermission = errors.New("无权限")
)

type NoteService struct {
	r *repository.NoteRepository
}

func NewNoteService(r *repository.NoteRepository) *NoteService {
	return &NoteService{r: r}
}

func (s *NoteService) Create(ctx context.Context, note *model.Note) error {
	return s.r.Create(ctx, note)
}

func (s *NoteService) CreateWithImages(ctx context.Context, note *model.Note) error {
	return s.r.CreateWithImages(ctx, note)
}

func (s *NoteService) DeleteByNoteId(ctx context.Context, noteId uint) error {
	return s.r.DeleteByNoteId(ctx, noteId)
}

func (s *NoteService) Update(ctx context.Context, note *model.Note) error {
	return s.r.Update(ctx, note)
}

func (s *NoteService) UpdateWithImages(ctx context.Context, note *model.Note) error {
	existingNote, err := s.r.GetById(ctx, note.ID)
	if err != nil {
		return ErrNoteNotFound
	}

	if existingNote.UserID != note.UserID {
		return ErrNoPermission
	}

	return s.r.UpdateWithImages(ctx, note)
}

func (s *NoteService) GetByUserId(ctx context.Context, userId uint) ([]model.Note, error) {
	return s.r.GetByUserId(ctx, userId)
}

func (s *NoteService) GetByUserIdWithPagination(ctx context.Context, userId uint, page, pageSize int) ([]model.Note, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return s.r.GetByUserIdWithPagination(ctx, userId, page, pageSize)
}

func (s *NoteService) GetLikedNotesByUserIdWithPagination(ctx context.Context, userId uint, page, pageSize int) ([]model.Note, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return s.r.GetLikedNotesByUserIdWithPagination(ctx, userId, page, pageSize)
}

func (s *NoteService) GetById(ctx context.Context, noteId uint) (*model.Note, error) {
	return s.r.GetById(ctx, noteId)
}

// Pagination 分页信息
type Pagination struct {
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	HasMore  bool  `json:"has_more"`
}

// NoteListResp 笔记列表响应结构体
type NoteListResp struct {
	List       []model.Note `json:"list"`
	Pagination Pagination   `json:"pagination"`
}

func (s *NoteService) SearchNotes(ctx context.Context, keyword string, page, pageSize int) (*NoteListResp, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	notes, total, err := s.r.SearchNotes(ctx, keyword, page, pageSize)
	if err != nil {
		return nil, err
	}

	totalPage := total / int64(pageSize)
	if total%int64(pageSize) != 0 {
		totalPage++
	}

	return &NoteListResp{
		List: notes,
		Pagination: Pagination{
			Total:    total,
			Page:     page,
			PageSize: pageSize,
			HasMore:  page < int(totalPage),
		},
	}, nil
}

func (s *NoteService) GetNoteList(ctx context.Context, page, pageSize int) (*NoteListResp, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	notes, total, err := s.r.GetNoteList(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	totalPage := total / int64(pageSize)
	if total%int64(pageSize) != 0 {
		totalPage++
	}

	return &NoteListResp{
		List: notes,
		Pagination: Pagination{
			Total:    total,
			Page:     page,
			PageSize: pageSize,
			HasMore:  page < int(totalPage),
		},
	}, nil
}
