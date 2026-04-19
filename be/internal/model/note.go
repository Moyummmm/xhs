package model

import (
	"gorm.io/gorm"
	"time"
)

type Note struct {
	// 覆盖 gorm.Model 中的字段，使用正确的 json 标签名
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Title        string `gorm:"size:100" json:"title"`
	Body         string `gorm:"size:500" json:"content"`
	CoverURL     string `gorm:"size:500" json:"cover_url"`
	VideoURL     string `gorm:"size:500" json:"video_url"`
	Location     string `gorm:"size:100" json:"location"`
	TopicID      uint   `gorm:"index;default:0" json:"topic_id"`
	Status       uint   `gorm:"default:1" json:"status"`
	UserID       uint   `gorm:"index" json:"user_id"`
	LikeCount    uint   `gorm:"default:0" json:"like_count"`
	CollectCount uint   `gorm:"default:0" json:"collect_count"`
	CommentCount uint   `gorm:"default:0" json:"comment_count"`

	// 关联User，不会再数据库创建外键约束
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	// 关联图片（一对多
	Images []NoteImage `gorm:"foreignKey:NoteID" json:"images,omitempty"`
}
