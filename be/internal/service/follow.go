package service

import (
	"context"
	"server/internal/model"
	"server/internal/repository"
)

type FollowService struct {
	r          *repository.FollowRepository
	userRepo   *repository.UserRepository
}

func NewFollowService(r *repository.FollowRepository, userRepo *repository.UserRepository) *FollowService {
	return &FollowService{r: r, userRepo: userRepo}
}

func (s *FollowService) Follow(ctx context.Context, followerId, followingId uint) error {
	if err := s.r.Follow(ctx, followerId, followingId); err != nil {
		return err
	}
	_ = s.userRepo.UpdateFollowerCount(ctx, followingId, 1)
	_ = s.userRepo.UpdateFollowingCount(ctx, followerId, 1)
	return nil
}

func (s *FollowService) Unfollow(ctx context.Context, followerId, followingId uint) error {
	if err := s.r.Unfollow(ctx, followerId, followingId); err != nil {
		return err
	}
	_ = s.userRepo.UpdateFollowerCount(ctx, followingId, -1)
	_ = s.userRepo.UpdateFollowingCount(ctx, followerId, -1)
	return nil
}

func (s *FollowService) GetFollowers(ctx context.Context, userId uint) ([]model.User, error) {
	return s.r.GetFollowers(ctx, userId)
}

func (s *FollowService) GetFollowings(ctx context.Context, userId uint) ([]model.User, error) {
	return s.r.GetFollowings(ctx, userId)
}

func (s *FollowService) GetFollowerCount(ctx context.Context, userId uint) (int, error) {
	user, err := s.userRepo.GetById(ctx, userId)
	if err != nil {
		return 0, err
	}
	return int(user.FollowerCount), nil
}

func (s *FollowService) GetFollowingCount(ctx context.Context, userId uint) (int, error) {
	user, err := s.userRepo.GetById(ctx, userId)
	if err != nil {
		return 0, err
	}
	return int(user.FollowingCount), nil
}
