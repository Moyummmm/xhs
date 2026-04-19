package repository

import (
	"server/internal/model"

	"gorm.io/gorm"
)

type FollowRepository struct {
	db *gorm.DB
}

func NewFollowRepository(db *gorm.DB) *FollowRepository {
	return &FollowRepository{db: db}
}

// Follow 关注用户
func (r *FollowRepository) Follow(followerId, followingId uint) error {
	follow := model.Follow{
		FollowerID:  followerId,
		FollowingID: followingId,
	}
	return r.db.Save(&follow).Error
}

// Unfollow 取消关注
func (r *FollowRepository) Unfollow(followerId, followingId uint) error {
	return r.db.Where("follower_id = ? AND following_id = ?", followerId, followingId).
		Delete(&model.Follow{}).Error
}

// GetFollowers 获取粉丝列表
func (r *FollowRepository) GetFollowers(userId uint) ([]model.User, error) {
	var users []model.User
	err := r.db.Joins("JOIN follows ON follows.follower_id = users.id").
		Where("follows.following_id = ?", userId).
		Find(&users).Error
	return users, err
}

// GetFollowings 获取关注列表
func (r *FollowRepository) GetFollowings(userId uint) ([]model.User, error) {
	var users []model.User
	err := r.db.Joins("JOIN follows ON follows.following_id = users.id").
		Where("follows.follower_id = ?", userId).
		Find(&users).Error
	return users, err
}

// GetFollowerCount 获取粉丝数量
func (r *FollowRepository) GetFollowerCount(userId uint) (int, error) {
	var count int64
	err := r.db.Model(&model.Follow{}).Where("following_id = ?", userId).Count(&count).Error
	return int(count), err
}

// GetFollowingCount 获取关注数量
func (r *FollowRepository) GetFollowingCount(userId uint) (int, error) {
	var count int64
	err := r.db.Model(&model.Follow{}).Where("follower_id = ?", userId).Count(&count).Error
	return int(count), err
}
