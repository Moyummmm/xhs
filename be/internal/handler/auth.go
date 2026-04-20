package handler

import (
	"server/config"
	"server/internal/repository"
	"server/internal/service"
	"server/pkg/errorConfig"
	"server/pkg/response"

	"github.com/gin-gonic/gin"
)

var authService *service.AuthService

func init() {
	userRepo := repository.NewUserRepository(config.DB())
	authService = service.NewAuthService(userRepo)
}

type RegisterReq struct {
	Nickname string `json:"nickname" binding:"required,min=3,max=30"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	ctx := c.Request.Context()

	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errorConfig.ErrBadRequest.Code, "请求参数错误")
		return
	}

	resp, err := authService.Register(ctx, service.RegisterReq{
		Username: req.Nickname,
		Password: req.Password,
	})
	if err != nil {
		code, msg := errorConfig.ExtractCodeAndMessage(err)
		response.Fail(c, code, msg)
		return
	}

	response.Success(c, resp)
}

func Login(c *gin.Context) {
	ctx := c.Request.Context()

	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errorConfig.ErrBadRequest.Code, "请求参数错误")
		return
	}

	resp, err := authService.Login(ctx, service.LoginReq{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		code, msg := errorConfig.ExtractCodeAndMessage(err)
		response.Fail(c, code, msg)
		return
	}

	response.Success(c, resp)
}

func Logout(c *gin.Context) {
	response.SuccessMsg(c, "登出成功", nil)
}

func RefreshToken(c *gin.Context) {
	ctx := c.Request.Context()

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		response.Fail(c, errorConfig.ErrUnauthorized.Code, "缺少Authorization头")
		return
	}

	var tokenString string
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = authHeader[7:]
	} else {
		tokenString = authHeader
	}

	resp, err := authService.RefreshToken(ctx, tokenString)
	if err != nil {
		code, msg := errorConfig.ExtractCodeAndMessage(err)
		response.Fail(c, code, msg)
		return
	}

	response.Success(c, gin.H{"token": resp.Token})
}
