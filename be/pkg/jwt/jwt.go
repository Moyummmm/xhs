package jwt

import (
	"fmt"
	"server/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

func InitJwt() error {
	jwtSecret = []byte(config.GlobalConfig.JWT.SecretKey)
	return nil
}

type Claims struct {
	UserID int64 `json:"userId"`
	jwt.RegisteredClaims
}

func GenerateToken(userId int64) (string, error) {
	expireTime, err := time.ParseDuration(config.GlobalConfig.JWT.ExpireTime)
	if err != nil {
		return "", fmt.Errorf("token expire time is invalid")
	}

	claims := &Claims{
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("token parse error: %s", err.Error())
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
