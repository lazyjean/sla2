package dto

import (
	"time"
)

// ResetPasswordRequest 重置密码请求
type ResetPasswordRequest struct {
	// 重置方式：phone-手机号验证码，apple-苹果登录
	ResetType string `json:"reset_type"`
	// 新密码
	NewPassword string `json:"new_password"`
	// 手机号（当reset_type为phone时必填）
	Phone string `json:"phone,omitempty"`
	// 验证码（当reset_type为phone时必填）
	VerificationCode string `json:"verification_code,omitempty"`
	// 苹果登录票据（当reset_type为apple时必填）
	AppleToken string `json:"apple_token,omitempty"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Nickname string `json:"nickname,omitempty"`
}

type RegisterResponse struct {
	UserID       uint32 `json:"user_id"`
	Token        string `json:"token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Account  string `json:"account"`  // 账号（用户名或邮箱）
	Password string `json:"password"` // 密码
}

// AuthResponse 认证响应
type LoginResponse struct {
	UserID        uint32 `json:"user_id"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Nickname      string `json:"nickname,omitempty"`
	Avatar        string `json:"avatar,omitempty"`
	Token         string `json:"token,omitempty"`
	RefreshToken  string `json:"refresh_token,omitempty"`
	IsNewUser     bool   `json:"is_new_user"`
}

// UpdateUserRequest 更新用户信息请求
type UpdateUserRequest struct {
	Nickname string `json:"nickname,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
}

type UpdateUserResponse struct {
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type ChangePasswordResponse struct {
}

type LogoutResponse struct {
}

type GetUserInfoRequest struct {
	UserID uint32 `json:"user_id"`
}

type GetUserInfoResponse struct {
	UserID        uint32 `json:"user_id"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Nickname      string `json:"nickname,omitempty"`
	Avatar        string `json:"avatar,omitempty"`
}

// RefreshTokenRequest 刷新token请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshTokenResponse 刷新token响应
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// UserDTO 用户数据传输对象
type UserDTO struct {
	ID            uint32    `json:"id"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	EmailVerified bool      `json:"email_verified"`
	Nickname      string    `json:"nickname,omitempty"`
	Avatar        string    `json:"avatar,omitempty"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// AppleLoginRequest Apple 登录请求
type AppleLoginRequest struct {
	AuthorizationCode string `json:"authorization_code"` // Apple Authorization Code
	UserIdentifier    string `json:"user_identifier"`    // Apple 用户唯一标识符
}

// AppleLoginResponse Apple 登录响应
type AppleLoginResponse struct {
	UserID       uint32 `json:"user_id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Nickname     string `json:"nickname,omitempty"`
	Avatar       string `json:"avatar,omitempty"`
	Token        string `json:"token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IsNewUser    bool   `json:"is_new_user"`
}
