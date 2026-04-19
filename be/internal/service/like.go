package service

import "server/internal/repository"

type LikeService struct {
	r *repository.LikeRepository
}

func NewLikeService(r *repository.LikeRepository) *LikeService {
	return &LikeService{r: r}
}

func (s *LikeService) LikeNote(userId, noteId uint) error {
	return s.r.LikeNote(userId, noteId)
}

func (s *LikeService) UnlikeNote(userId, noteId uint) error {
	return s.r.UnlikeNote(userId, noteId)
}

func (s *LikeService) GetLikeCountByNoteId(noteId uint) (int, error) {
	return s.r.GetLikeCountByNoteId(noteId)
}

func (s *LikeService) IsLiked(userId, noteId uint) (bool, error) {
	return s.r.IsLiked(userId, noteId)
}
