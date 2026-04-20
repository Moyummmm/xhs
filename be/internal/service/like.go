package service

import (
	"context"
	"server/internal/repository"
)

type LikeService struct {
	r *repository.LikeRepository
}

func NewLikeService(r *repository.LikeRepository) *LikeService {
	return &LikeService{r: r}
}

func (s *LikeService) LikeNote(ctx context.Context, userId, noteId uint) error {
	return s.r.LikeNote(ctx, userId, noteId)
}

func (s *LikeService) UnlikeNote(ctx context.Context, userId, noteId uint) error {
	return s.r.UnlikeNote(ctx, userId, noteId)
}

func (s *LikeService) GetLikeCountByNoteId(ctx context.Context, noteId uint) (int, error) {
	return s.r.GetLikeCountByNoteId(ctx, noteId)
}

func (s *LikeService) IsLiked(ctx context.Context, userId, noteId uint) (bool, error) {
	return s.r.IsLiked(ctx, userId, noteId)
}
