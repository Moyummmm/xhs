package service

import (
	"server/internal/model"
	"server/internal/repository"
)

type FollowService struct {
	r *repository.FollowRepository
}

func NewFollowService(r *repository.FollowRepository) *FollowService {
	return &FollowService{r: r}
}

func (s *FollowService) Follow(followerId, followingId uint) error {
	return s.r.Follow(followerId, followingId)
}

func (s *FollowService) Unfollow(followerId, followingId uint) error {
	return s.r.Unfollow(followerId, followingId)
}

func (s *FollowService) GetFollowers(userId uint) ([]model.User, error) {
	return s.r.GetFollowers(userId)
}

func (s *FollowService) GetFollowings(userId uint) ([]model.User, error) {
	return s.r.GetFollowings(userId)
}

func (s *FollowService) GetFollowerCount(userId uint) (int, error) {
	return s.r.GetFollowerCount(userId)
}

func (s *FollowService) GetFollowingCount(userId uint) (int, error) {
	return s.r.GetFollowingCount(userId)
}
