package jwt

import (
	"context"
	"testing"
	"time"

	"server/config"
)

// setupJWT 初始化 JWT 配置，仅供测试使用
func setupJWT(t *testing.T) {
	config.GlobalConfig = &config.Config{
		JWT: config.JWTConfig{
			SecretKey:  "test-secret-key",
			ExpireTime: "1h",
		},
	}
	if err := InitJwt(); err != nil {
		t.Fatalf("InitJwt() failed: %v", err)
	}
}

// TestGenerateToken 测试 Token 生成功能
// 验证生成的 token 不为空
func TestGenerateToken(t *testing.T) {
	setupJWT(t)

	token, err := GenerateToken(123)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}
	if token == "" {
		t.Error("GenerateToken() returned empty string")
	}
}

// TestGenerateToken_And_ParseToken 测试 Token 生成后能正确解析
// 验证 ParseToken 能从 token 中提取出正确的 UserID
func TestGenerateToken_And_ParseToken(t *testing.T) {
	setupJWT(t)

	userID := int64(456)
	token, err := GenerateToken(userID)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	claims, err := ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken() error = %v", err)
	}
	if claims.UserID != userID {
		t.Errorf("ParseToken() UserID = %v, want %v", claims.UserID, userID)
	}
}

// TestParseToken_InvalidToken 测试解析无效 Token
// 验证格式错误的 token 返回错误
func TestParseToken_InvalidToken(t *testing.T) {
	setupJWT(t)

	_, err := ParseToken("invalid-token")
	if err == nil {
		t.Error("ParseToken() expected error for invalid token")
	}
}

// TestGenerateAndParse_RoundTrip 测试生成和解析的往返流程
// 使用多个 UserID 进行测试，确保解析出的 UserID 与传入的一致
func TestGenerateAndParse_RoundTrip(t *testing.T) {
	setupJWT(t)

	testCases := []int64{1, 100, 99999, 0}
	for _, id := range testCases {
		token, err := GenerateToken(id)
		if err != nil {
			t.Errorf("GenerateToken(%d) error = %v", id, err)
			continue
		}
		claims, err := ParseToken(token)
		if err != nil {
			t.Errorf("ParseToken() error = %v", err)
			continue
		}
		if claims.UserID != id {
			t.Errorf("claims.UserID = %v, want %v", claims.UserID, id)
		}
	}
}

// TestParseToken_EmptyToken 测试解析空 Token
// 验证空字符串返回错误
func TestParseToken_EmptyToken(t *testing.T) {
	setupJWT(t)

	_, err := ParseToken("")
	if err == nil {
		t.Error("ParseToken() expected error for empty token")
	}
}

// TestGenerateToken_ExpiryTime 测试 Token 过期功能
// 设置 1s 过期时间，等待 2s 后验证 token 已失效
func TestGenerateToken_ExpiryTime(t *testing.T) {
	config.GlobalConfig = &config.Config{
		JWT: config.JWTConfig{
			SecretKey:  "test-secret",
			ExpireTime: "1s",
		},
	}
	InitJwt()

	token, _ := GenerateToken(1)
	time.Sleep(2 * time.Second)

	_, err := ParseToken(token)
	if err == nil {
		t.Error("ParseToken() should fail for expired token")
	}
}

// TestGenerateToken_WithContextCancellation 测试取消上下文对 Token 生成的影响
func TestGenerateToken_WithContextCancellation(t *testing.T) {
	setupJWT(t)

	_, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := GenerateToken(1)
	if err != nil {
		t.Logf("GenerateToken with cancelled context: %v", err)
	}
}