package password

import (
	"testing"
)

// TestHash 测试密码哈希生成功能
// 覆盖场景：正常密码、空字符串密码、短密码
func TestHash(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"正常密码", "password123", false},
		{"空字符串", "", true},
		{"短密码", "123", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := Hash(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Hash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && hash == "" {
				t.Error("Hash() returned empty string")
			}
		})
	}
}

// TestHash_Uniqueness 测试相同密码生成的哈希值不同
// bcrypt 会自动加盐，相同输入应产生不同哈希
func TestHash_Uniqueness(t *testing.T) {
	hash1, _ := Hash("password")
	hash2, _ := Hash("password")
	if hash1 == hash2 {
		t.Error("Hash() should produce different hashes for same input (bcrypt adds salt)")
	}
}

// TestVerify 测试密码验证功能
// 覆盖场景：正确密码、错误密码、空密码、空hash
func TestVerify(t *testing.T) {
	hash, _ := Hash("password123")

	tests := []struct {
		name     string
		password string
		hashed   string
		want     bool
	}{
		{"正确密码", "password123", hash, true},
		{"错误密码", "wrongpassword", hash, false},
		{"空密码", "", hash, false},
		{"空hash", "password123", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Verify(tt.password, tt.hashed); got != tt.want {
				t.Errorf("Verify() = %v, want %v", got, tt.want)
			}
		})
	}
}