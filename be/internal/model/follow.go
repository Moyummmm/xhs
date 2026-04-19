package model

import "gorm.io/gorm"

type Follow struct {
	gorm.Model
	// 关注者
	FollowerID uint `gorm:"index;uniqueIndex:idx_follower_following;not null" json:"follower_id"`
	// 被关注着
	FollowingID uint `gorm:"index;uniqueIndex:idx_follower_following;not null" json:"following_id"`
}
