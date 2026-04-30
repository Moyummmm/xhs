package service

import (
	"context"
	"server/internal/cache"
	"server/internal/model"
	"server/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetById(ctx context.Context, id uint) (*model.User, error) {
	// Try cache first
	if cached, err := cache.GetUser(ctx, id); err == nil && cached != nil {
		return cached.ToModel(), nil
	}
	// Cache miss: query DB
	user, err := s.userRepo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	// Store in cache
	cache.SetUser(ctx, id, cache.NewCachedUser(user))
	return user, nil
}

func (s *UserService) Patch(ctx context.Context, user model.User) (bool, error) {
	result, err := s.userRepo.PatchByUsername(ctx, user)
	if result {
		cache.DeleteUser(ctx, user.ID)
	}
	return result, err
}

func (s *UserService) UpdateById(ctx context.Context, id uint, user model.User) (*model.User, error) {
	updated, err := s.userRepo.UpdateById(ctx, id, user)
	if err == nil {
		cache.DeleteUser(ctx, id)
	}
	return updated, err
}

func (s *UserService) DeleteById(ctx context.Context, id uint) error {
	err := s.userRepo.DeleteById(ctx, id)
	if err == nil {
		cache.DeleteUser(ctx, id)
	}
	return err
}
