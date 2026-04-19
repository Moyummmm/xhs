package errorConfig

import (
	"fmt"
)

// AppError 自定义业务错误
// 用于统一服务层、控制层的错误返回，包含业务状态码和用户提示信息
type AppError struct {
	Code    int    // 业务错误码（非 HTTP 状态码）
	Message string // 给用户看的错误提示
	Err     error  // 底层原始错误（可选，用于日志记录和排查问题）
}

// Error 实现 error 接口
// 让 AppError 可以被当作标准 error 类型使用
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[code:%d] %s | cause: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[code:%d] %s", e.Code, e.Message)
}

// Unwrap 返回底层错误
// 支持 errors.Is 和 errors.As 链式判断
func (e *AppError) Unwrap() error {
	return e.Err
}

// WithMessage 复制当前错误并修改提示信息
// 用于在传递过程中补充更具体的上下文
func (e *AppError) WithMessage(msg string) *AppError {
	return &AppError{
		Code:    e.Code,
		Message: msg,
		Err:     e.Err,
	}
}

// WithError 复制当前错误并包装底层错误
// 用于在传递过程中附加原始错误
func (e *AppError) WithError(err error) *AppError {
	return &AppError{
		Code:    e.Code,
		Message: e.Message,
		Err:     err,
	}
}

// ============================================
// 预定义常见业务错误
// 控制层可直接使用这些错误返回统一格式的 JSON 响应
// ============================================

var (
	// ErrInternalServer 服务器内部错误
	ErrInternalServer = New(10001, "服务器繁忙，请稍后再试")

	// ErrBadRequest 请求参数错误
	ErrBadRequest = New(10002, "请求参数错误")

	// ErrUnauthorized 未登录或 token 无效
	ErrUnauthorized = New(10003, "请先登录")

	// ErrForbidden 无权限访问
	ErrForbidden = New(10004, "无权限访问")

	// ErrNotFound 资源不存在
	ErrNotFound = New(10005, "请求的资源不存在")

	// ErrTooManyRequests 请求过于频繁
	ErrTooManyRequests = New(10006, "操作过于频繁，请稍后再试")

	// ErrDuplicate 数据重复（如重复点赞、重复关注）
	ErrDuplicate = New(10007, "重复操作")

	// ErrDatabase 数据库操作失败
	ErrDatabase = New(10008, "数据操作失败")

	// ErrInvalidPassword 密码错误
	ErrInvalidPassword = New(10009, "账号或密码错误")

	// ErrUserNotExist 用户不存在
	ErrUserNotExist = New(10010, "用户不存在")
)

// New 创建一个自定义业务错误
// code: 业务错误码
// message: 用户提示信息
func New(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Wrap 包装一个已有的标准 error 为业务错误
// 示例：errors.Wrap(ErrDatabase, sqlErr)
func Wrap(appErr *AppError, err error) *AppError {
	return appErr.WithError(err)
}

// IsAppError 判断是否是自定义业务错误
func IsAppError(err error) (*AppError, bool) {
	if err == nil {
		return nil, false
	}
	if e, ok := err.(*AppError); ok {
		return e, true
	}
	return nil, false
}

// GetCode 从 error 中提取业务码
// 如果不是 AppError，返回默认内部错误码 10001
func GetCode(err error) int {
	if e, ok := IsAppError(err); ok {
		return e.Code
	}
	return ErrInternalServer.Code
}

// GetMessage 从 error 中提取用户提示信息
// 如果不是 AppError，返回默认的友好提示
func GetMessage(err error) string {
	if e, ok := IsAppError(err); ok {
		return e.Message
	}
	return ErrInternalServer.Message
}

// ExtractCodeAndMessage 从 error 中提取 code 和 message
// 方便 handler 层统一处理
func ExtractCodeAndMessage(err error) (int, string) {
	if e, ok := IsAppError(err); ok {
		return e.Code, e.Message
	}
	return ErrInternalServer.Code, ErrInternalServer.Message
}
