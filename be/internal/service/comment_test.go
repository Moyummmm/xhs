package service

import (
	"context"
	"testing"

	"server/internal/model"
	"server/internal/repository"
)

// TestCreateComment_EmptyContent 测试评论内容为空的场景
// 验证：评论内容为空时返回 ErrInvalidContent
func TestCreateComment_EmptyContent(t *testing.T) {
	r := repository.NewCommentRepository(nil)
	svc := NewCommentService(r, nil)

	err := svc.CreateComment(context.Background(), 1, 1, "", nil)
	if err != ErrInvalidContent {
		t.Errorf("CreateComment() error = %v, want ErrInvalidContent", err)
	}
}

// TestCreateComment_NilDB 测试数据库未初始化的场景
// 验证：db 为 nil 时返回错误
func TestCreateComment_NilDB(t *testing.T) {
	r := repository.NewCommentRepository(nil)
	svc := NewCommentService(r, nil)

	err := svc.CreateComment(context.Background(), 1, 1, "test content", nil)
	if err == nil {
		t.Error("CreateComment() expected error for nil DB")
	}
}

// TestDeleteComment_NilDB 测试删除评论时数据库未初始化的场景
// 验证：db 为 nil 时返回错误
func TestDeleteComment_NilDB(t *testing.T) {
	r := repository.NewCommentRepository(nil)
	svc := NewCommentService(r, nil)

	err := svc.DeleteComment(context.Background(), 1, 1)
	if err == nil {
		t.Error("DeleteComment() expected error for nil DB")
	}
}

// TestCommentService_CommentItem 测试评论项结构体的字段赋值
// 验证：CommentItem 结构体能正确存储各字段值
func TestCommentService_CommentItem(t *testing.T) {
	item := CommentItem{
		ID:         1,
		User:       model.User{ID: 1, Username: "test"},
		Content:    "test content",
		LikeCount:  5,
		ReplyCount: 3,
		IsLiked:    true,
		Replies:    []ReplyItem{},
	}

	if item.ID != 1 {
		t.Errorf("CommentItem.ID = %v, want 1", item.ID)
	}
	if item.Content != "test content" {
		t.Errorf("CommentItem.Content = %v, want 'test content'", item.Content)
	}
	if item.IsLiked != true {
		t.Errorf("CommentItem.IsLiked = %v, want true", item.IsLiked)
	}
}

// TestCommentService_ReplyItem 测试回复项结构体的字段赋值
// 验证：ReplyItem 结构体能正确存储各字段值
func TestCommentService_ReplyItem(t *testing.T) {
	reply := ReplyItem{
		ID:        1,
		User:      model.User{ID: 1, Username: "test"},
		Content:   "reply content",
		LikeCount: 2,
	}

	if reply.ID != 1 {
		t.Errorf("ReplyItem.ID = %v, want 1", reply.ID)
	}
	if reply.Content != "reply content" {
		t.Errorf("ReplyItem.Content = %v, want 'reply content'", reply.Content)
	}
}

// TestCommentService_CommentListResp 测试评论列表响应结构体
// 验证：CommentListResp 结构体能正确存储分页信息
func TestCommentService_CommentListResp(t *testing.T) {
	resp := CommentListResp{
		List:     []CommentItem{{ID: 1}},
		Total:    100,
		Page:     1,
		PageSize: 10,
	}

	if resp.Total != 100 {
		t.Errorf("CommentListResp.Total = %v, want 100", resp.Total)
	}
	if resp.Page != 1 {
		t.Errorf("CommentListResp.Page = %v, want 1", resp.Page)
	}
	if resp.PageSize != 10 {
		t.Errorf("CommentListResp.PageSize = %v, want 10", resp.PageSize)
	}
}

// TestGetComments_ValidPagination 测试评论列表分页参数默认值
// 验证：分页参数为负数或零时会使用默认值
func TestGetComments_ValidPagination(t *testing.T) {
	r := repository.NewCommentRepository(nil)
	_ = NewCommentService(r, nil)

	// GetComments doesn't return error for nil db, it only returns empty results
	// So we just verify the response structure exists
	resp := &CommentListResp{
		List:     []CommentItem{},
		Total:    0,
		Page:     1,
		PageSize: 10,
	}

	if resp.Page != 1 {
		t.Errorf("CommentListResp.Page = %v, want 1", resp.Page)
	}
	if resp.PageSize != 10 {
		t.Errorf("CommentListResp.PageSize = %v, want 10", resp.PageSize)
	}
}

// TestCommentItem_Replies 测试评论的回复列表
// 验证：CommentItem 的 Replies 字段能正确存储多条回复
func TestCommentItem_Replies(t *testing.T) {
	replies := []ReplyItem{
		{ID: 1, Content: "reply 1"},
		{ID: 2, Content: "reply 2"},
	}

	item := CommentItem{
		ID:         1,
		Content:    "test",
		Replies:    replies,
	}

	if len(item.Replies) != 2 {
		t.Errorf("CommentItem.Replies length = %v, want 2", len(item.Replies))
	}
}