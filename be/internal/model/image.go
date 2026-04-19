package model

import (
	"gorm.io/gorm"
	"time"
)

type Image struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	URL    string `gorm:"size:500;not null" json:"url"`
	Width  int    `gorm:"default:0" json:"width"`
	Height int    `gorm:"default:0" json:"height"`
	UserID uint   `gorm:"index" json:"user_id"`
}
