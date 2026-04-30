package repository

import (
	"context"
	"server/internal/model"
)

// UserRepositoryInterface 定义 UserRepository 的接口，用于依赖注入和测试
type UserRepositoryInterface interface {
	Create(ctx context.Context, user *model.User) error
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	PatchByUsername(ctx context.Context, user model.User) (bool, error)
	GetById(ctx context.Context, id uint) (*model.User, error)
	DeleteById(ctx context.Context, id uint) error
	UpdateById(ctx context.Context, id uint, user model.User) (*model.User, error)
	UpdateFollowerCount(ctx context.Context, userId uint, delta int) error
	UpdateFollowingCount(ctx context.Context, userId uint, delta int) error
}

// 确保 UserRepository 实现 UserRepositoryInterface
var _ UserRepositoryInterface = (*UserRepository)(nil)