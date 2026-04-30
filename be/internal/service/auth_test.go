package service

import (
	"context"
	"errors"
	"testing"

	"server/config"
	"server/internal/model"
	"server/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository 模拟 repository.UserRepositoryInterface，用于单元测试
type MockUserRepository struct {
	CreateFunc           func(ctx context.Context, user *model.User) error
	GetByUsernameFunc    func(ctx context.Context, username string) (*model.User, error)
	ExistsByUsernameFunc func(ctx context.Context, username string) (bool, error)
	GetByIdFunc          func(ctx context.Context, id uint) (*model.User, error)
	DeleteByIdFunc       func(ctx context.Context, id uint) error
	PatchByUsernameFunc  func(ctx context.Context, user model.User) (bool, error)
	UpdateByIdFunc       func(ctx context.Context, id uint, user model.User) (*model.User, error)
	UpdateFollowerCountFunc  func(ctx context.Context, userId uint, delta int) error
	UpdateFollowingCountFunc func(ctx context.Context, userId uint, delta int) error
}

// Create 实现 UserRepositoryInterface 的 Create 方法
func (m *MockUserRepository) Create(ctx context.Context, user *model.User) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, user)
	}
	return nil
}

// GetByUsername 实现 UserRepositoryInterface 的 GetByUsername 方法
func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	if m.GetByUsernameFunc != nil {
		return m.GetByUsernameFunc(ctx, username)
	}
	return nil, errors.New("not implemented")
}

// ExistsByUsername 实现 UserRepositoryInterface 的 ExistsByUsername 方法
func (m *MockUserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	if m.ExistsByUsernameFunc != nil {
		return m.ExistsByUsernameFunc(ctx, username)
	}
	return false, nil
}

// GetById 实现 UserRepositoryInterface 的 GetById 方法
func (m *MockUserRepository) GetById(ctx context.Context, id uint) (*model.User, error) {
	if m.GetByIdFunc != nil {
		return m.GetByIdFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

// DeleteById 实现 UserRepositoryInterface 的 DeleteById 方法
func (m *MockUserRepository) DeleteById(ctx context.Context, id uint) error {
	return nil
}

// PatchByUsername 实现 UserRepositoryInterface 的 PatchByUsername 方法
func (m *MockUserRepository) PatchByUsername(ctx context.Context, user model.User) (bool, error) {
	return true, nil
}

// UpdateById 实现 UserRepositoryInterface 的 UpdateById 方法
func (m *MockUserRepository) UpdateById(ctx context.Context, id uint, user model.User) (*model.User, error) {
	return &user, nil
}

// UpdateFollowerCount 实现 UserRepositoryInterface 的 UpdateFollowerCount 方法
func (m *MockUserRepository) UpdateFollowerCount(ctx context.Context, userId uint, delta int) error {
	return nil
}

// UpdateFollowingCount 实现 UserRepositoryInterface 的 UpdateFollowingCount 方法
func (m *MockUserRepository) UpdateFollowingCount(ctx context.Context, userId uint, delta int) error {
	return nil
}

// Interface compliance check - 确保 MockUserRepository 实现了 UserRepositoryInterface
var _ interface {
	Create(ctx context.Context, user *model.User) error
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	GetById(ctx context.Context, id uint) (*model.User, error)
	DeleteById(ctx context.Context, id uint) error
	PatchByUsername(ctx context.Context, user model.User) (bool, error)
	UpdateById(ctx context.Context, id uint, user model.User) (*model.User, error)
	UpdateFollowerCount(ctx context.Context, userId uint, delta int) error
	UpdateFollowingCount(ctx context.Context, userId uint, delta int) error
} = (*MockUserRepository)(nil)

// initJWT 初始化 JWT 配置，供测试使用
func initJWT(t *testing.T) {
	config.GlobalConfig = &config.Config{
		JWT: config.JWTConfig{
			SecretKey:  "test-secret-key",
			ExpireTime: "1h",
		},
	}
	if err := jwt.InitJwt(); err != nil {
		t.Fatalf("InitJwt() failed: %v", err)
	}
}

// TestAuthService_Register_Success 测试成功注册场景
// 验证：用户名不存在时能成功创建用户并返回 token
func TestAuthService_Register_Success(t *testing.T) {
	initJWT(t)

	mockRepo := &MockUserRepository{
		ExistsByUsernameFunc: func(ctx context.Context, username string) (bool, error) {
			return false, nil
		},
		CreateFunc: func(ctx context.Context, user *model.User) error {
			user.ID = 1
			return nil
		},
		GetByIdFunc: func(ctx context.Context, id uint) (*model.User, error) {
			return &model.User{ID: id, Username: "testuser"}, nil
		},
	}

	svc := NewAuthService(mockRepo)
	resp, err := svc.Register(context.Background(), RegisterReq{
		Username: "testuser",
		Password: "password123",
	})

	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	if resp == nil {
		t.Fatal("Register() resp is nil")
	}
	if resp.Token == "" {
		t.Error("Register() token is empty")
	}
	if resp.User.ID != 1 {
		t.Errorf("Register() User.ID = %v, want 1", resp.User.ID)
	}
	if resp.User.Username != "testuser" {
		t.Errorf("Register() User.Username = %v, want testuser", resp.User.Username)
	}
}

// TestAuthService_Register_UsernameExists 测试用户名已存在场景
// 验证：用户名已存在时返回错误
func TestAuthService_Register_UsernameExists(t *testing.T) {
	initJWT(t)

	mockRepo := &MockUserRepository{
		ExistsByUsernameFunc: func(ctx context.Context, username string) (bool, error) {
			return true, nil
		},
	}

	svc := NewAuthService(mockRepo)
	_, err := svc.Register(context.Background(), RegisterReq{
		Username: "existinguser",
		Password: "password123",
	})

	if err == nil {
		t.Fatal("Register() expected error for existing username")
	}
}

// TestAuthService_Login_Success 测试成功登录场景
// 验证：用户名和密码正确时返回 token 和用户信息
func TestAuthService_Login_Success(t *testing.T) {
	initJWT(t)

	hashedPassword, _ := hashPasswordForTest("password123")
	mockRepo := &MockUserRepository{
		GetByUsernameFunc: func(ctx context.Context, username string) (*model.User, error) {
			if username == "testuser" {
				return &model.User{
					ID:       1,
					Username: "testuser",
					Password: hashedPassword,
				}, nil
			}
			return nil, errors.New("user not found")
		},
		GetByIdFunc: func(ctx context.Context, id uint) (*model.User, error) {
			return &model.User{ID: id, Username: "testuser"}, nil
		},
	}

	svc := NewAuthService(mockRepo)
	resp, err := svc.Login(context.Background(), LoginReq{
		Username: "testuser",
		Password: "password123",
	})

	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if resp == nil {
		t.Fatal("Login() resp is nil")
	}
	if resp.Token == "" {
		t.Error("Login() token is empty")
	}
	if resp.User.Username != "testuser" {
		t.Errorf("Login() User.Username = %v, want testuser", resp.User.Username)
	}
}

// TestAuthService_Login_WrongPassword 测试密码错误场景
// 验证：密码错误时返回错误
func TestAuthService_Login_WrongPassword(t *testing.T) {
	initJWT(t)

	hashedPassword, _ := hashPasswordForTest("correctpassword")
	mockRepo := &MockUserRepository{
		GetByUsernameFunc: func(ctx context.Context, username string) (*model.User, error) {
			return &model.User{
				ID:       1,
				Username: "testuser",
				Password: hashedPassword,
			}, nil
		},
	}

	svc := NewAuthService(mockRepo)
	_, err := svc.Login(context.Background(), LoginReq{
		Username: "testuser",
		Password: "wrongpassword",
	})

	if err == nil {
		t.Fatal("Login() expected error for wrong password")
	}
}

// TestAuthService_Login_UserNotFound 测试用户不存在场景
// 验证：用户不存在时返回错误
func TestAuthService_Login_UserNotFound(t *testing.T) {
	initJWT(t)

	mockRepo := &MockUserRepository{
		GetByUsernameFunc: func(ctx context.Context, username string) (*model.User, error) {
			return nil, errors.New("user not found")
		},
	}

	svc := NewAuthService(mockRepo)
	_, err := svc.Login(context.Background(), LoginReq{
		Username: "nonexistent",
		Password: "password123",
	})

	if err == nil {
		t.Fatal("Login() expected error for non-existent user")
	}
}

// TestAuthService_RefreshToken_Success 测试成功刷新 Token 场景
// 验证：有效的 token 能成功刷新并返回新的 token
func TestAuthService_RefreshToken_Success(t *testing.T) {
	initJWT(t)

	user := &model.User{ID: 1, Username: "testuser"}
	mockRepo := &MockUserRepository{
		GetByIdFunc: func(ctx context.Context, id uint) (*model.User, error) {
			return user, nil
		},
	}

	svc := NewAuthService(mockRepo)
	token, _ := jwt.GenerateToken(1)
	resp, err := svc.RefreshToken(context.Background(), token)

	if err != nil {
		t.Fatalf("RefreshToken() error = %v", err)
	}
	if resp == nil {
		t.Fatal("RefreshToken() resp is nil")
	}
	if resp.Token == "" {
		t.Error("RefreshToken() token is empty")
	}
}

// TestAuthService_RefreshToken_InvalidToken 测试无效 Token 场景
// 验证：无效的 token 返回错误
func TestAuthService_RefreshToken_InvalidToken(t *testing.T) {
	initJWT(t)

	mockRepo := &MockUserRepository{}
	svc := NewAuthService(mockRepo)

	_, err := svc.RefreshToken(context.Background(), "invalid-token")

	if err == nil {
		t.Fatal("RefreshToken() expected error for invalid token")
	}
}

// TestAuthService_RefreshToken_UserNotFound 测试 Token 对应用户不存在场景
// 验证：token 有效但用户已被删除时返回错误
func TestAuthService_RefreshToken_UserNotFound(t *testing.T) {
	initJWT(t)

	mockRepo := &MockUserRepository{
		GetByIdFunc: func(ctx context.Context, id uint) (*model.User, error) {
			return nil, errors.New("user not found")
		},
	}

	svc := NewAuthService(mockRepo)
	token, _ := jwt.GenerateToken(999)
	_, err := svc.RefreshToken(context.Background(), token)

	if err == nil {
		t.Fatal("RefreshToken() expected error for non-existent user")
	}
}

// hashPasswordForTest 测试辅助函数，用于生成 bcrypt 哈希密码
func hashPasswordForTest(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}