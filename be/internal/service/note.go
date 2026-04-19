package service

import (
	"errors"
	"server/internal/model"
	"server/internal/repository"
)

var (
	ErrNoteNotFound = errors.New("笔记不存在")
	ErrNoPermission = errors.New("无权限")
)

// NoteService 笔记服务结构体，封装笔记相关的业务逻辑
type NoteService struct {
	r *repository.NoteRepository // 笔记仓库实例，用于数据持久化操作
}

// NewNoteService 创建一个新的笔记服务实例
// 参数:
//   - r: 笔记仓库实例
//
// 返回:
//   - 笔记服务实例指针
func NewNoteService(r *repository.NoteRepository) *NoteService {
	return &NoteService{r: r}
}

// Create 创建新笔记
// 参数:
//   - note: 笔记模型指针，包含笔记的详细信息
//
// 返回:
//   - 错误信息，如果创建成功则返回nil
func (s *NoteService) Create(note *model.Note) error {
	return s.r.Create(note)
}

// CreateWithImages 创建笔记并附带图片
func (s *NoteService) CreateWithImages(note *model.Note) error {
	return s.r.CreateWithImages(note)
}

// DeleteByNoteId 根据笔记ID删除笔记
// 参数:
//   - noteId: 笔记的唯一标识ID
//
// 返回:
//   - 错误信息，如果删除成功则返回nil
func (s *NoteService) DeleteByNoteId(noteId uint) error {
	return s.r.DeleteByNoteId(noteId)
}

// Update 更新笔记信息
// 参数:
//   - note: 笔记模型指针，包含更新后的笔记信息
//
// 返回:
//   - 错误信息，如果更新成功则返回nil
func (s *NoteService) Update(note *model.Note) error {
	return s.r.Update(note)
}

// UpdateWithImages 更新笔记并附带图片
func (s *NoteService) UpdateWithImages(note *model.Note) error {
	existingNote, err := s.r.GetById(note.ID)
	if err != nil {
		return ErrNoteNotFound
	}

	if existingNote.UserID != note.UserID {
		return ErrNoPermission
	}

	return s.r.UpdateWithImages(note)
}

// GetByUserId 根据用户ID获取该用户的所有笔记
// 参数:
//   - userId: 用户的唯一标识ID
//
// 返回:
//   - 笔记列表，包含该用户的所有笔记
//   - 错误信息，如果查询成功则返回nil
func (s *NoteService) GetByUserId(userId uint) ([]model.Note, error) {
	return s.r.GetByUserId(userId)
}

// GetByUserIdWithPagination 根据用户ID分页获取笔记
func (s *NoteService) GetByUserIdWithPagination(userId uint, page, pageSize int) ([]model.Note, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return s.r.GetByUserIdWithPagination(userId, page, pageSize)
}

// GetById 根据笔记ID获取笔记
func (s *NoteService) GetById(noteId uint) (*model.Note, error) {
	return s.r.GetById(noteId)
}

// Pagination 分页信息
type Pagination struct {
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	HasMore  bool  `json:"has_more"`
}

// NoteListResp 笔记列表响应结构体
type NoteListResp struct {
	List       []model.Note `json:"list"`
	Pagination Pagination   `json:"pagination"`
}

// SearchNotes 搜索笔记
func (s *NoteService) SearchNotes(keyword string, page, pageSize int) (*NoteListResp, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	notes, total, err := s.r.SearchNotes(keyword, page, pageSize)
	if err != nil {
		return nil, err
	}

	totalPage := total / int64(pageSize)
	if total%int64(pageSize) != 0 {
		totalPage++
	}

	return &NoteListResp{
		List: notes,
		Pagination: Pagination{
			Total:    total,
			Page:     page,
			PageSize: pageSize,
			HasMore:  page < int(totalPage),
		},
	}, nil
}

// GetNoteList 获取笔记列表，按综合打分排序
func (s *NoteService) GetNoteList(page, pageSize int) (*NoteListResp, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	notes, total, err := s.r.GetNoteList(page, pageSize)
	if err != nil {
		return nil, err
	}

	totalPage := total / int64(pageSize)
	if total%int64(pageSize) != 0 {
		totalPage++
	}

	return &NoteListResp{
		List: notes,
		Pagination: Pagination{
			Total:    total,
			Page:     page,
			PageSize: pageSize,
			HasMore:  page < int(totalPage),
		},
	}, nil
}
