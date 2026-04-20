package repository

import (
	"context"
	"server/internal/model"

	"gorm.io/gorm"
)

type NoteRepository struct {
	db *gorm.DB
}

func NewNoteRepository(db *gorm.DB) *NoteRepository {
	return &NoteRepository{db: db}
}

func (r *NoteRepository) Create(ctx context.Context, note *model.Note) error {
	return r.db.WithContext(ctx).Create(note).Error
}

func (r *NoteRepository) CreateWithImages(ctx context.Context, note *model.Note) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(note).Error; err != nil {
			return err
		}
		for i := range note.Images {
			note.Images[i].NoteID = note.ID
		}
		if len(note.Images) > 0 {
			if err := tx.Create(&note.Images).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *NoteRepository) Update(ctx context.Context, note *model.Note) error {
	return r.db.WithContext(ctx).Updates(note).Error
}

func (r *NoteRepository) Delete(ctx context.Context, note *model.Note) error {
	return r.db.WithContext(ctx).Delete(note).Error
}

func (r *NoteRepository) GetByUserId(ctx context.Context, userid uint) ([]model.Note, error) {
	var notes []model.Note
	err := r.db.WithContext(ctx).Where("user_id = ?", userid).Find(&notes).Error
	return notes, err
}

func (r *NoteRepository) GetByUserIdWithPagination(ctx context.Context, userid uint, page, pageSize int) ([]model.Note, int64, error) {
	var notes []model.Note
	var total int64

	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Model(&model.Note{}).Where("user_id = ?", userid).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.WithContext(ctx).Where("user_id = ?", userid).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&notes).Error

	return notes, total, err
}

func (r *NoteRepository) GetLikedNotesByUserIdWithPagination(ctx context.Context, userId uint, page, pageSize int) ([]model.Note, int64, error) {
	var notes []model.Note
	var total int64

	offset := (page - 1) * pageSize

	subQuery := r.db.WithContext(ctx).Model(&model.Like4Note{}).Select("note_id").Where("user_id = ?", userId)

	if err := r.db.WithContext(ctx).Model(&model.Note{}).Where("id IN (?)", subQuery).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.WithContext(ctx).Preload("User").Preload("Images").
		Where("id IN (?)", subQuery).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&notes).Error

	return notes, total, err
}

func (r *NoteRepository) GetById(ctx context.Context, noteId uint) (*model.Note, error) {
	var note model.Note
	err := r.db.WithContext(ctx).Preload("User").Preload("Images").Where("id = ?", noteId).First(&note).Error
	if err != nil {
		return nil, err
	}
	return &note, nil
}

func (r *NoteRepository) DeleteByNoteId(ctx context.Context, noteId uint) error {
	return r.db.WithContext(ctx).Where("id = ?", noteId).Delete(&model.Note{}).Error
}

func (r *NoteRepository) UpdateWithImages(ctx context.Context, note *model.Note) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Note{}).Where("id = ?", note.ID).Updates(map[string]interface{}{
			"title":     note.Title,
			"body":      note.Body,
			"cover_url": note.CoverURL,
			"video_url": note.VideoURL,
			"location":  note.Location,
			"topic_id":  note.TopicID,
			"user_id":   note.UserID,
		}).Error; err != nil {
			return err
		}

		if err := tx.Where("note_id = ?", note.ID).Delete(&model.NoteImage{}).Error; err != nil {
			return err
		}

		if len(note.Images) > 0 {
			for i := range note.Images {
				note.Images[i].NoteID = note.ID
			}
			if err := tx.Create(&note.Images).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *NoteRepository) SearchNotes(ctx context.Context, keyword string, page, pageSize int) ([]model.Note, int64, error) {
	var notes []model.Note
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.WithContext(ctx).Model(&model.Note{}).Where("title LIKE ? OR body LIKE ?", "%"+keyword+"%", "%"+keyword+"%")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&notes).Error
	return notes, total, err
}

func (r *NoteRepository) GetNoteList(ctx context.Context, page, pageSize int) ([]model.Note, int64, error) {
	var notes []model.Note
	var total int64

	offset := (page - 1) * pageSize

	if err := r.db.WithContext(ctx).Model(&model.Note{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.WithContext(ctx).Preload("User").Preload("Images").
		Order("like_count * 3 + collect_count * 5 DESC, created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&notes).Error

	return notes, total, err
}
