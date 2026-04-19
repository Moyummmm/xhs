package model

import "gorm.io/gorm"

type Collect struct {
	gorm.Model

	Active bool `gorm:"default:true" json:"active"`
	// uniqueIndex:idx_user_note_collect 确保collect的幂等性
	NoteID uint `gorm:"index;uniqueIndex:idx_user_note_collect;not null" json:"note_id"`
	UserID uint `gorm:"index;uniqueIndex:idx_user_note_collect;not null" json:"user_id"`
}
