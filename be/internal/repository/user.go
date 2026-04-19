package repository

import (
	"server/internal/model"

	"gorm.io/gorm"
)

// UserRepository 用户仓库
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓库
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create 创建用户
func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// GetByUsername 根据用户名查询用户
func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// ExistsByUsername 检查用户名是否存在
func (r *UserRepository) ExistsByUsername(username string) (bool, error) {
	var count int64
	if err := r.db.Model(&model.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// PatchByUsername 依据用户名来修改信息
func (r *UserRepository) PatchByUsername(user model.User) (bool, error) {
	result := r.db.Model(&model.User{}).Where("username = ?", user.Username).Updates(user)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

// GetById 根据ID查询用户
func (r *UserRepository) GetById(id uint) (*model.User, error) {
	var user model.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// DeleteById 根据ID删除用户
func (r *UserRepository) DeleteById(id uint) error {
	return r.db.Where("id = ?", id).Delete(&model.User{}).Error
}

// UpdateById 根据ID更新用户信息
func (r *UserRepository) UpdateById(id uint, user model.User) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Updates(user).Error
}
