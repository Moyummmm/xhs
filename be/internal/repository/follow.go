package repository

import (
	"context"
	"server/internal/model"

	"gorm.io/gorm"
)

type FollowRepository struct {
	db *gorm.DB
}

func NewFollowRepository(db *gorm.DB) *FollowRepository {
	return &FollowRepository{db: db}
}

func (r *FollowRepository) Follow(ctx context.Context, followerId, followingId uint) error {
	follow := model.Follow{
		FollowerID:  followerId,
		FollowingID: followingId,
	}
	return r.db.WithContext(ctx).Save(&follow).Error
}

func (r *FollowRepository) Unfollow(ctx context.Context, followerId, followingId uint) error {
	return r.db.WithContext(ctx).Where("follower_id = ? AND following_id = ?", followerId, followingId).
		Delete(&model.Follow{}).Error
}

func (r *FollowRepository) GetFollowers(ctx context.Context, userId uint) ([]model.User, error) {
	var users []model.User
	err := r.db.WithContext(ctx).Joins("JOIN follows ON follows.follower_id = users.id").
		Where("follows.following_id = ?", userId).
		Find(&users).Error
	return users, err
}

func (r *FollowRepository) GetFollowings(ctx context.Context, userId uint) ([]model.User, error) {
	var users []model.User
	err := r.db.WithContext(ctx).Joins("JOIN follows ON follows.following_id = users.id").
		Where("follows.follower_id = ?", userId).
		Find(&users).Error
	return users, err
}

func (r *FollowRepository) GetFollowerCount(ctx context.Context, userId uint) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Follow{}).Where("following_id = ?", userId).Count(&count).Error
	return int(count), err
}

func (r *FollowRepository) GetFollowingCount(ctx context.Context, userId uint) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Follow{}).Where("follower_id = ?", userId).Count(&count).Error
	return int(count), err
}
