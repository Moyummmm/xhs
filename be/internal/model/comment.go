package model

import (
	"gorm.io/gorm"
	"time"
)

type Comment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Content  string `gorm:"size:500;not null" json:"content"`
	NoteID   uint   `gorm:"index;not null" json:"note_id"`
	UserID   uint   `gorm:"index;not null" json:"user_id"`
	ParentID *uint  `gorm:"index;default:0" json:"parent_id"`

	// 【新增】评论的点赞数
	LikeCount uint `gorm:"default:0" json:"like_count"`
	// 回复数（一级评论特有）
	ReplyCount uint `gorm:"default:0" json:"reply_count"`

	User   User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Parent *Comment `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
}
