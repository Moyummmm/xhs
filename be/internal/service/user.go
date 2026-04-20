package service

import (
	"context"
	"server/internal/model"
	"server/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func (s *UserService) Patch(ctx context.Context, user model.User) (bool, error) {
	return s.userRepo.PatchByUsername(ctx, user)
}

func (s *UserService) GetById(ctx context.Context, id uint) (*model.User, error) {
	return s.userRepo.GetById(ctx, id)
}

func (s *UserService) UpdateById(ctx context.Context, id uint, user model.User) (*model.User, error) {
	return s.userRepo.UpdateById(ctx, id, user)
}

func (s *UserService) DeleteById(ctx context.Context, id uint) error {
	return s.userRepo.DeleteById(ctx, id)
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}
