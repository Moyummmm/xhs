package repository

import (
	"server/internal/model"

	"gorm.io/gorm"
)

// NoteRepository 笔记仓库，封装对笔记数据的数据库操作
type NoteRepository struct {
	db *gorm.DB
}

// NewNoteRepository 创建一个新的笔记仓库实例
func NewNoteRepository(db *gorm.DB) *NoteRepository {
	return &NoteRepository{db: db}
}

// Create 创建一条新的笔记记录
func (r *NoteRepository) Create(note *model.Note) error {
	return r.db.Create(note).Error
}

// CreateWithImages 创建笔记并附带图片（使用事务）
func (r *NoteRepository) CreateWithImages(note *model.Note) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
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

// Update 更新笔记记录
// 注意：Updates 方法只会更新非零值字段，如果需要更新零值字段，应使用 Save 方法
func (r *NoteRepository) Update(note *model.Note) error {
	return r.db.Updates(note).Error
}

// Delete 删除笔记记录
func (r *NoteRepository) Delete(note *model.Note) error {
	return r.db.Delete(note).Error
}

// GetByUserId 根据用户ID获取该用户的所有笔记
func (r *NoteRepository) GetByUserId(userid uint) ([]model.Note, error) {
	var notes []model.Note
	err := r.db.Where("user_id = ?", userid).Find(&notes).Error
	return notes, err
}

// GetByUserIdWithPagination 根据用户ID分页获取笔记
func (r *NoteRepository) GetByUserIdWithPagination(userid uint, page, pageSize int) ([]model.Note, int64, error) {
	var notes []model.Note
	var total int64

	offset := (page - 1) * pageSize
	if err := r.db.Model(&model.Note{}).Where("user_id = ?", userid).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Where("user_id = ?", userid).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&notes).Error

	return notes, total, err
}

// GetById 根据笔记ID获取笔记
func (r *NoteRepository) GetById(noteId uint) (*model.Note, error) {
	var note model.Note
	err := r.db.Preload("User").Preload("Images").Where("id = ?", noteId).First(&note).Error
	if err != nil {
		return nil, err
	}
	return &note, nil
}

// DeleteByNoteId 根据笔记ID删除笔记记录
func (r *NoteRepository) DeleteByNoteId(noteId uint) error {
	return r.db.Where("id = ?", noteId).Delete(&model.Note{}).Error
}

// UpdateWithImages 更新笔记并附带图片（使用事务）
func (r *NoteRepository) UpdateWithImages(note *model.Note) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
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

// SearchNotes 搜索笔记
func (r *NoteRepository) SearchNotes(keyword string, page, pageSize int) ([]model.Note, int64, error) {
	var notes []model.Note
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&model.Note{}).Where("title LIKE ? OR body LIKE ?", "%"+keyword+"%", "%"+keyword+"%")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&notes).Error
	return notes, total, err
}

// GetNoteList 获取笔记列表，按综合打分排序
// 综合打分公式: score = like_count * 3 + collect_count * 5 + 时间衰减分数
func (r *NoteRepository) GetNoteList(page, pageSize int) ([]model.Note, int64, error) {
	var notes []model.Note
	var total int64

	offset := (page - 1) * pageSize

	if err := r.db.Model(&model.Note{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Preload("User").Preload("Images").
		Order("like_count * 3 + collect_count * 5 DESC, created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&notes).Error

	return notes, total, err
}
