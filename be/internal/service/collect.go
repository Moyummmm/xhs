package service

import (
	"server/internal/model"
	"server/internal/repository"
)

type CollectService struct {
	r *repository.CollectRepository
}

func NewCollectService(r *repository.CollectRepository) *CollectService {
	return &CollectService{r: r}
}

func (s *CollectService) CollectById(userId, noteid uint) error {
	return s.r.CollectById(userId, noteid)
}

func (s *CollectService) DisCollectById(userId, noteId uint) error {
	return s.r.DisCollectById(userId, noteId)
}

func (s *CollectService) GetCollectedCount(userId uint) (int, error) {
	return s.r.GetCollecCountsByUserId(userId)
}

func (s *CollectService) GetCollectList(userId uint) ([]model.Note, error) {
	return s.r.GetCollectList(userId)
}
