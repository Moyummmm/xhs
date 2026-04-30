package service

import (
	"context"
	"server/internal/model"
	"server/internal/repository"
	"server/pkg/errorConfig"
	"server/pkg/jwt"
	"server/pkg/password"
)

type AuthService struct {
	userRepo repository.UserRepositoryInterface
}

func NewAuthService(userRepo repository.UserRepositoryInterface) *AuthService {
	return &AuthService{userRepo: userRepo}
}

type RegisterReq struct {
	Username string
	Password string
}

type RegisterResp struct {
	Token string     `json:"token"`
	User  model.User `json:"user"`
}

func (s *AuthService) Register(ctx context.Context, req RegisterReq) (*RegisterResp, error) {
	exist, err := s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, errorConfig.ErrDatabase
	}
	if exist {
		return nil, errorConfig.ErrBadRequest.WithMessage("用户名已存在")
	}

	hash, err := password.Hash(req.Password)
	if err != nil {
		return nil, errorConfig.ErrInternalServer
	}

	user := model.User{
		Username: req.Username,
		Password: hash,
	}
	if err := s.userRepo.Create(ctx, &user); err != nil {
		return nil, errorConfig.ErrDatabase
	}

	token, err := jwt.GenerateToken(int64(user.ID))
	if err != nil {
		return nil, errorConfig.ErrInternalServer
	}

	return &RegisterResp{
		Token: token,
		User:  user,
	}, nil
}

type LoginReq struct {
	Username string
	Password string
}

type LoginResp struct {
	Token string     `json:"token"`
	User  model.User `json:"user"`
}

func (s *AuthService) Login(ctx context.Context, req LoginReq) (*LoginResp, error) {
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, errorConfig.ErrInvalidPassword
	}

	if !password.Verify(req.Password, user.Password) {
		return nil, errorConfig.ErrInvalidPassword
	}

	token, err := jwt.GenerateToken(int64(user.ID))
	if err != nil {
		return nil, errorConfig.ErrInternalServer
	}

	return &LoginResp{
		Token: token,
		User:  *user,
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, tokenString string) (*LoginResp, error) {
	claims, err := jwt.ParseToken(tokenString)
	if err != nil {
		return nil, errorConfig.ErrUnauthorized.WithMessage("无效的token")
	}

	user, err := s.userRepo.GetById(ctx, uint(claims.UserID))
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