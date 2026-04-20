package service

import (
	"context"
	"server/internal/model"
	"server/internal/repository"
)

type FollowService struct {
	r *repository.FollowRepository
}

func NewFollowService(r *repository.FollowRepository) *FollowService {
	return &FollowService{r: r}
}

func (s *FollowService) Follow(ctx context.Context, followerId, followingId uint) error {
	return s.r.Follow(ctx, followerId, followingId)
}

func (s *FollowService) Unfollow(ctx context.Context, followerId, followingId uint) error {
	return s.r.Unfollow(ctx, followerId, followingId)
}

func (s *FollowService) GetFollowers(ctx context.Context, userId uint) ([]model.User, error) {
	return s.r.GetFollowers(ctx, userId)
}

func (s *FollowService) GetFollowings(ctx context.Context, userId uint) ([]model.User, error) {
	return s.r.GetFollowings(ctx, userId)
}

func (s *FollowService) GetFollowerCount(ctx context.Context, userId uint) (int, error) {
	return s.r.GetFollowerCount(ctx, userId)
}

func (s *FollowService) GetFollowingCount(ctx context.Context, userId uint) (int, error) {
	return s.r.GetFollowingCount(ctx, userId)
}
