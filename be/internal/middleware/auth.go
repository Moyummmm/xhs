package middleware

import (
	"net/http"
	"strings"

	"server/pkg/jwt"

	"github.com/gin-gonic/gin"
)

const ContextUserIDKey = "userId"

// Auth 用户认证中间件
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractToken(c)
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": "401",
				"msg":  "login first",
			})
			return
		}

		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "登录过期",
			})
			return
		}
		c.Set(ContextUserIDKey, claims.UserID)
		c.Next()
	}
}

func extractToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
			return parts[1]
		}
		return authHeader
	}

	if token := c.GetHeader("X-Token"); token != "" {
		return token
	}

	if token := c.GetHeader("AccessToken"); token != "" {
		return token
	}
	return c.Query("token")
}

func CurrentUserID(c *gin.Context) int64 {
	val, exists := c.Get(ContextUserIDKey)
	if !exists {
		return 0
	}
	uid, ok := val.(int64)
	if !ok {
		return 0
	}
	return uid
}
