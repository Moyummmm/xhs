package service

import (
	"context"
	"server/internal/model"
	"server/internal/repository"
)

type CollectService struct {
	r *repository.CollectRepository
}

func NewCollectService(r *repository.CollectRepository) *CollectService {
	return &CollectService{r: r}
}

func (s *CollectService) CollectById(ctx context.Context, userId, noteid uint) error {
	return s.r.CollectById(ctx, userId, noteid)
}

func (s *CollectService) DisCollectById(ctx context.Context, userId, noteId uint) error {
	return s.r.DisCollectById(ctx, userId, noteId)
}

func (s *CollectService) GetCollectedCount(ctx context.Context, userId uint) (int, error) {
	return s.r.GetCollecCountsByUserId(ctx, userId)
}

func (s *CollectService) GetCollectList(ctx context.Context, userId uint) ([]model.Note, error) {
	return s.r.GetCollectList(ctx, userId)
}

func (s *CollectService) GetCollectListWithPagination(ctx context.Context, userId uint, page, pageSize int) ([]model.Note, int64, error) {
	return s.r.GetCollectListWithPagination(ctx, userId, page, pageSize)
}
