package service

import (
	"server/internal/model"
	"server/internal/repository"
	"server/pkg/errorConfig"
	"server/pkg/jwt"
	"server/pkg/password"
)

// AuthService 认证服务
type AuthService struct {
	userRepo *repository.UserRepository
}

// NewAuthService 创建认证服务
func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

// RegisterReq 注册请求
type RegisterReq struct {
	Username string
	Password string
}

// RegisterResp 注册响应
type RegisterResp struct {
	Token string     `json:"token"`
	User  model.User `json:"user"`
}

// Register 用户注册
func (s *AuthService) Register(req RegisterReq) (*RegisterResp, error) {
	// 1. 检查用户名是否已存在
	exist, err := s.userRepo.ExistsByUsername(req.Username)
	if err != nil {
		return nil, errorConfig.ErrDatabase
	}
	if exist {
		return nil, errorConfig.ErrBadRequest.WithMessage("用户名已存在")
	}

	// 2. 密码哈希
	hash, err := password.Hash(req.Password)
	if err != nil {
		return nil, errorConfig.ErrInternalServer
	}

	// 3. 创建用户
	user := model.User{
		Username: req.Username,
		Password: hash,
	}
	if err := s.userRepo.Create(&user); err != nil {
		return nil, errorConfig.ErrDatabase
	}

	// 4. 生成 JWT token
	token, err := jwt.GenerateToken(int64(user.ID))
	if err != nil {
		return nil, errorConfig.ErrInternalServer
	}

	return &RegisterResp{
		Token: token,
		User:  user,
	}, nil
}

// LoginReq 登录请求
type LoginReq struct {
	Username string
	Password string
}

// LoginResp 登录响应
type LoginResp struct {
	Token string     `json:"token"`
	User  model.User `json:"user"`
}

// Login 用户登录
func (s *AuthService) Login(req LoginReq) (*LoginResp, error) {
	// 1. 根据用户名查询用户
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		return nil, errorConfig.ErrInvalidPassword
	}

	// 2. 验证密码
	if !password.Verify(req.Password, user.Password) {
		return nil, errorConfig.ErrInvalidPassword
	}

	// 3. 生成 JWT token
	token, err := jwt.GenerateToken(int64(user.ID))
	if err != nil {
		return nil, errorConfig.ErrInternalServer
	}

	return &LoginResp{
		Token: token,
		User:  *user,
	}, nil
}

// RefreshToken 刷新token
func (s *AuthService) RefreshToken(tokenString string) (*LoginResp, error) {
	claims, err := jwt.ParseToken(tokenString)
	if err != nil {
		return nil, errorConfig.ErrUnauthorized.WithMessage("无效的token")
	}

	user, err := s.userRepo.GetById(uint(claims.UserID))
	if err != nil {
		return nil, errorConfig.ErrNotFound.WithMessage("用户不存在")
	}

	newToken, err := jwt.GenerateToken(claims.UserID)
	if err != nil {
		return nil, errorConfig.ErrInternalServer
	}

	return &LoginResp{
		Token: newToken,
		User:  *user,
	}, nil
}
