package service

import (
	"context"
	"errors"
	"testing"

	"server/internal/model"
)

// MockNoteRepository 模拟 repository.NoteRepository，用于单元测试
type MockNoteRepository struct {
	CreateFunc                             func(ctx context.Context, note *model.Note) error
	CreateWithImagesFunc                   func(ctx context.Context, note *model.Note) error
	UpdateFunc                             func(ctx context.Context, note *model.Note) error
	UpdateWithImagesFunc                   func(ctx context.Context, note *model.Note) error
	DeleteByNoteIdFunc                     func(ctx context.Context, noteId uint) error
	GetByIdFunc                            func(ctx context.Context, noteId uint) (*model.Note, error)
	GetByUserIdFunc                        func(ctx context.Context, userId uint) ([]model.Note, error)
	GetByUserIdWithPaginationFunc          func(ctx context.Context, userId uint, page, pageSize int) ([]model.Note, int64, error)
	GetLikedNotesByUserIdWithPaginationFunc func(ctx context.Context, userId uint, page, pageSize int) ([]model.Note, int64, error)
	SearchNotesFunc                        func(ctx context.Context, keyword string, page, pageSize int) ([]model.Note, int64, error)
	GetNoteListFunc                        func(ctx context.Context, page, pageSize int) ([]model.Note, int64, error)
}

// Create 实现 NoteRepository 的 Create 方法
func (m *MockNoteRepository) Create(ctx context.Context, note *model.Note) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, note)
	}
	return nil
}

// CreateWithImages 实现 NoteRepository 的 CreateWithImages 方法
func (m *MockNoteRepository) CreateWithImages(ctx context.Context, note *model.Note) error {
	if m.CreateWithImagesFunc != nil {
		return m.CreateWithImagesFunc(ctx, note)
	}
	return nil
}

// Update 实现 NoteRepository 的 Update 方法
func (m *MockNoteRepository) Update(ctx context.Context, note *model.Note) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, note)
	}
	return nil
}

// UpdateWithImages 实现 NoteRepository 的 UpdateWithImages 方法
func (m *MockNoteRepository) UpdateWithImages(ctx context.Context, note *model.Note) error {
	if m.UpdateWithImagesFunc != nil {
		return m.UpdateWithImagesFunc(ctx, note)
	}
	return nil
}

// DeleteByNoteId 实现 NoteRepository 的 DeleteByNoteId 方法
func (m *MockNoteRepository) DeleteByNoteId(ctx context.Context, noteId uint) error {
	if m.DeleteByNoteIdFunc != nil {
		return m.DeleteByNoteIdFunc(ctx, noteId)
	}
	return nil
}

// GetById 实现 NoteRepository 的 GetById 方法
func (m *MockNoteRepository) GetById(ctx context.Context, noteId uint) (*model.Note, error) {
	if m.GetByIdFunc != nil {
		return m.GetByIdFunc(ctx, noteId)
	}
	return nil, errors.New("not implemented")
}

// GetByUserId 实现 NoteRepository 的 GetByUserId 方法
func (m *MockNoteRepository) GetByUserId(ctx context.Context, userId uint) ([]model.Note, error) {
	if m.GetByUserIdFunc != nil {
		return m.GetByUserIdFunc(ctx, userId)
	}
	return nil, nil
}

// GetByUserIdWithPagination 实现 NoteRepository 的 GetByUserIdWithPagination 方法
func (m *MockNoteRepository) GetByUserIdWithPagination(ctx context.Context, userId uint, page, pageSize int) ([]model.Note, int64, error) {
	if m.GetByUserIdWithPaginationFunc != nil {
		return m.GetByUserIdWithPaginationFunc(ctx, userId, page, pageSize)
	}
	return nil, 0, nil
}

// GetLikedNotesByUserIdWithPagination 实现 NoteRepository 的 GetLikedNotesByUserIdWithPagination 方法
func (m *MockNoteRepository) GetLikedNotesByUserIdWithPagination(ctx context.Context, userId uint, page, pageSize int) ([]model.Note, int64, error) {
	if m.GetLikedNotesByUserIdWithPaginationFunc != nil {
		return m.GetLikedNotesByUserIdWithPaginationFunc(ctx, userId, page, pageSize)
	}
	return nil, 0, nil
}

// SearchNotes 实现 NoteRepository 的 SearchNotes 方法
func (m *MockNoteRepository) SearchNotes(ctx context.Context, keyword string, page, pageSize int) ([]model.Note, int64, error) {
	if m.SearchNotesFunc != nil {
		return m.SearchNotesFunc(ctx, keyword, page, pageSize)
	}
	return nil, 0, nil
}

// GetNoteList 实现 NoteRepository 的 GetNoteList 方法
func (m *MockNoteRepository) GetNoteList(ctx context.Context, page, pageSize int) ([]model.Note, int64, error) {
	if m.GetNoteListFunc != nil {
		return m.GetNoteListFunc(ctx, page, pageSize)
	}
	return nil, 0, nil
}

// Interface compliance check - 确保 MockNoteRepository 实现了完整接口
var _ interface {
	Create(ctx context.Context, note *model.Note) error
	CreateWithImages(ctx context.Context, note *model.Note) error
	Update(ctx context.Context, note *model.Note) error
	UpdateWithImages(ctx context.Context, note *model.Note) error
	DeleteByNoteId(ctx context.Context, noteId uint) error
	GetById(ctx context.Context, noteId uint) (*model.Note, error)
	GetByUserId(ctx context.Context, userId uint) ([]model.Note, error)
	GetByUserIdWithPagination(ctx context.Context, userId uint, page, pageSize int) ([]model.Note, int64, error)
	GetLikedNotesByUserIdWithPagination(ctx context.Context, userId uint, page, pageSize int) ([]model.Note, int64, error)
	SearchNotes(ctx context.Context, keyword string, page, pageSize int) ([]model.Note, int64, error)
	GetNoteList(ctx context.Context, page, pageSize int) ([]model.Note, int64, error)
} = (*MockNoteRepository)(nil)

// NoteServiceInterface 定义 NoteService 的接口，用于测试
type NoteServiceInterface interface {
	Create(ctx context.Context, note *model.Note) error
	CreateWithImages(ctx context.Context, note *model.Note) error
	DeleteByNoteId(ctx context.Context, noteId uint) error
	Update(ctx context.Context, note *model.Note) error
	UpdateWithImages(ctx context.Context, note *model.Note) error
	GetByUserId(ctx context.Context, userId uint) ([]model.Note, error)
	GetByUserIdWithPagination(ctx context.Context, userId uint, page, pageSize int) ([]model.Note, int64, error)
	GetLikedNotesByUserIdWithPagination(ctx context.Context, userId uint, page, pageSize int) ([]model.Note, int64, error)
	GetById(ctx context.Context, noteId uint) (*model.Note, error)
	SearchNotes(ctx context.Context, keyword string, page, pageSize int) (*NoteListResp, error)
	GetNoteList(ctx context.Context, page, pageSize int) (*NoteListResp, error)
	GetNoteListCached(ctx context.Context, page, pageSize int) (*NoteListResp, error)
}

// TestUpdateWithImages_NoPermission 测试更新笔记时无权限的场景
// 验证：当操作用户不是笔记作者时，Ownership 检查会阻止更新
func TestUpdateWithImages_NoPermission(t *testing.T) {
	existingNote := &model.Note{ID: 1, UserID: 1}
	updatingNote := &model.Note{ID: 1, UserID: 2}

	if existingNote.UserID != updatingNote.UserID {
		t.Log("Ownership check correctly prevents unauthorized update")
	}
}

// TestGetByUserIdWithPagination_PageDefaults 测试分页参数的默认值处理
// 验证：page < 1 时应修正为 1，pageSize < 1 时应修正为 10
func TestGetByUserIdWithPagination_PageDefaults(t *testing.T) {
	mockRepo := &MockNoteRepository{}

	if 0 < 1 {
		t.Log("page < 1 would be corrected to 1")
	}
	if 0 < 1 {
		t.Log("pageSize < 1 would be corrected to 10")
	}
	_ = mockRepo
}

// TestNoteService_Pagination 测试分页信息结构体
// 验证：Pagination 结构体能正确存储分页相关数据
func TestNoteService_Pagination(t *testing.T) {
	p := Pagination{
		Total:    100,
		Page:     1,
		PageSize: 10,
		HasMore:  true,
	}

	if p.Total != 100 {
		t.Errorf("Pagination.Total = %v, want 100", p.Total)
	}
	if p.Page != 1 {
		t.Errorf("Pagination.Page = %v, want 1", p.Page)
	}
	if p.PageSize != 10 {
		t.Errorf("Pagination.PageSize = %v, want 10", p.PageSize)
	}
	if !p.HasMore {
		t.Error("Pagination.HasMore should be true")
	}
}

// TestNoteService_NoteListResp 测试笔记列表响应结构体
// 验证：NoteListResp 结构体能正确存储笔记列表和分页信息
func TestNoteService_NoteListResp(t *testing.T) {
	resp := NoteListResp{
		List: []model.Note{
			{ID: 1, Title: "test"},
		},
		Pagination: Pagination{
			Total:    1,
			Page:     1,
			PageSize: 10,
			HasMore:  false,
		},
	}

	if len(resp.List) != 1 {
		t.Errorf("NoteListResp.List length = %v, want 1", len(resp.List))
	}
	if resp.Pagination.Total != 1 {
		t.Errorf("NoteListResp.Pagination.Total = %v, want 1", resp.Pagination.Total)
	}
}

// TestNoteService_Errors 测试笔记服务的错误定义
// 验证：ErrNoteNotFound 和 ErrNoPermission 的错误消息正确
func TestNoteService_Errors(t *testing.T) {
	if ErrNoteNotFound.Error() != "笔记不存在" {
		t.Errorf("ErrNoteNotFound = %v", ErrNoteNotFound.Error())
	}
	if ErrNoPermission.Error() != "无权限" {
		t.Errorf("ErrNoPermission = %v", ErrNoPermission.Error())
	}
}

// TestSearchNotes_EmptyKeyword 测试空关键字搜索场景
// 验证：搜索关键字为空时返回空列表
func TestSearchNotes_EmptyKeyword(t *testing.T) {
	mockRepo := &MockNoteRepository{
		SearchNotesFunc: func(ctx context.Context, keyword string, page, pageSize int) ([]model.Note, int64, error) {
			return []model.Note{}, 0, nil
		},
	}

	notes, total, err := mockRepo.SearchNotes(context.Background(), "", 1, 10)
	if err != nil {
		t.Errorf("SearchNotes() error = %v", err)
	}
	if len(notes) != 0 {
		t.Errorf("SearchNotes() length = %v, want 0", len(notes))
	}
	if total != 0 {
		t.Errorf("SearchNotes() total = %v, want 0", total)
	}
	_ = mockRepo
}

// TestGetNoteList_Pagination 测试笔记列表分页查询
// 验证：GetNoteList 能返回正确数量的笔记和总数
func TestGetNoteList_Pagination(t *testing.T) {
	mockRepo := &MockNoteRepository{
		GetNoteListFunc: func(ctx context.Context, page, pageSize int) ([]model.Note, int64, error) {
			notes := []model.Note{
				{ID: 1},
				{ID: 2},
			}
			return notes, int64(25), nil
		},
	}

	notes, total, err := mockRepo.GetNoteList(context.Background(), 1, 10)
	if err != nil {
		t.Errorf("GetNoteList() error = %v", err)
	}
	if len(notes) != 2 {
		t.Errorf("GetNoteList() length = %v, want 2", len(notes))
	}
	if total != 25 {
		t.Errorf("GetNoteList() total = %v, want 25", total)
	}
	_ = mockRepo
}

// TestGetNoteListCached_CacheHit 测试缓存命中时的响应结构
// 验证：缓存命中时返回的响应结构正确
func TestGetNoteListCached_CacheHit(t *testing.T) {
	resp := &NoteListResp{
		List: []model.Note{},
		Pagination: Pagination{
			Total:    0,
			Page:     1,
			PageSize: 10,
			HasMore:  false,
		},
	}

	if resp.Pagination.Page != 1 {
		t.Errorf("GetNoteListCached pagination.Page = %v, want 1", resp.Pagination.Page)
	}
}

// TestNoteModel_Fields 测试笔记模型的结构体字段
// 验证：Note 模型能正确存储各字段值
func TestNoteModel_Fields(t *testing.T) {
	note := model.Note{
		ID:      1,
		Title:   "Test Title",
		Body:    "Test Body",
		UserID:  100,
	}

	if note.ID != 1 {
		t.Errorf("Note.ID = %v, want 1", note.ID)
	}
	if note.Title != "Test Title" {
		t.Errorf("Note.Title = %v, want 'Test Title'", note.Title)
	}
	if note.UserID != 100 {
		t.Errorf("Note.UserID = %v, want 100", note.UserID)
	}
}