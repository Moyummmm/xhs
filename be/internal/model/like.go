package model

import "gorm.io/gorm"

type Like4Note struct {
	gorm.Model
	NoteID uint `gorm:"uniqueIndex:idx_user_note;not null" json:"note_id"`
	UserID uint `gorm:"uniqueIndex:idx_user_note;not null" json:"user_id"`
}

type Like4Comment struct {
	gorm.Model
	CommentID uint `gorm:"uniqueIndex:idx_user_comment;not null" json:"comment_id"`
	UserID    uint `gorm:"uniqueIndex:idx_user_comment;not null" json:"user_id"`
}
