package utils

import "testing"

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  bool
	}{
		{
			name:  "有效的邮箱地址",
			email: "test@example.com",
			want:  true,
		},
		{
			name:  "有效的邮箱地址带点号",
			email: "test.user@example.com",
			want:  true,
		},
		{
			name:  "无效的邮箱地址-没有@",
			email: "testexample.com",
			want:  false,
		},
		{
			name:  "无效的邮箱地址-没有域名",
			email: "test@",
			want:  false,
		},
		{
			name:  "无效的邮箱地址-特殊字符",
			email: "test*@example.com",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidEmail(tt.email); got != tt.want {
				t.Errorf("IsValidEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidPhone(t *testing.T) {
	tests := []struct {
		name  string
		phone string
		want  bool
	}{
		{
			name:  "有效的手机号",
			phone: "13812345678",
			want:  true,
		},
		{
			name:  "有效的手机号-不同运营商",
			phone: "18912345678",
			want:  true,
		},
		{
			name:  "无效的手机号-位数不够",
			phone: "1381234567",
			want:  false,
		},
		{
			name:  "无效的手机号-位数过多",
			phone: "138123456789",
			want:  false,
		},
		{
			name:  "无效的手机号-非法前缀",
			phone: "12812345678",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidPhone(tt.phone); got != tt.want {
				t.Errorf("IsValidPhone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		want     bool
	}{
		{
			name:     "有效的用户名",
			username: "user123",
			want:     true,
		},
		{
			name:     "有效的用户名-带下划线",
			username: "user_123",
			want:     true,
		},
		{
			name:     "有效的用户名-带横线",
			username: "user-123",
			want:     true,
		},
		{
			name:     "无效的用户名-太短",
			username: "us",
			want:     false,
		},
		{
			name:     "无效的用户名-太长",
			username: "usernameiswaytoolongtobevalid",
			want:     false,
		},
		{
			name:     "无效的用户名-特殊字符",
			username: "user@123",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidUsername(tt.username); got != tt.want {
				t.Errorf("IsValidUsername() = %v, want %v", got, tt.want)
			}
		})
	}
}
