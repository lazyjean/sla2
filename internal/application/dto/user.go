package dto

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

// LoginRequest 登录请求
type LoginRequest struct {
	Account  string `json:"account"`  // 账号（用户名或邮箱）
	Password string `json:"password"` // 密码
}

// AuthResponse 认证响应
type AuthResponse struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Nickname string `json:"nickname,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
	Token    string `json:"token,omitempty"`
}

// UpdateUserRequest 更新用户信息请求
type UpdateUserRequest struct {
	Nickname string `json:"nickname,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
