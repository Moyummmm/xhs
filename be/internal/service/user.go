package service

import (
	"server/internal/model"
	"server/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

// Patch 修改信息
func (s *UserService) Patch(user model.User) (bool, error) {
	return s.userRepo.PatchByUsername(user)
}

// GetById 根据ID获取用户
func (s *UserService) GetById(id uint) (*model.User, error) {
	return s.userRepo.GetById(id)
}

// UpdateById 根据ID更新用户信息
func (s *UserService) UpdateById(id uint, user model.User) error {
	return s.userRepo.UpdateById(id, user)
}

// DeleteById 根据ID删除用户
func (s *UserService) DeleteById(id uint) error {
	return s.userRepo.DeleteById(id)
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}
