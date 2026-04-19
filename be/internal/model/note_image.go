package model

import (
	"gorm.io/gorm"
	"time"
)

type NoteImage struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	NoteID uint   `gorm:"index;not null" json:"note_id"`
	URL    string `gorm:"size:500" json:"url"`
	Width  int    `gorm:"default:800" json:"width"`
	Height int    `gorm:"default:600" json:"height"`
}