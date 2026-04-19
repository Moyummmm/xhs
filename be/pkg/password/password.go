package password

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Hash 使用 bcrypt 生成密码哈希
// 适合作密码存储，自动加盐
func Hash(password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("password is empty")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// Verify 验证密码
// 自动处理盐，使用常数时间比较防止时序攻击
func Verify(password, hashed string) bool {
	if password == "" || hashed == "" {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}
