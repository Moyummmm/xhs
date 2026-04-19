package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Username       string `gorm:"size:50;uniqueIndex" json:"nickname"`
	Password       string `gorm:"size:255" json:"-"`
	Avatar         string `gorm:"size:500" json:"avatar"`
	Bio            string `gorm:"size:500" json:"bio"`
	Gender         uint   `gorm:"default:0" json:"gender"`
	Birthday       string `gorm:"size:20" json:"birthday"`
	FollowingCount uint   `gorm:"default:0" json:"following_count"`
	FollowerCount  uint   `gorm:"default:0" json:"follower_count"`
}
