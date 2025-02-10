package dto

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string
	Password string
	Email    string
	Nickname string
}

// LoginRequest 登录请求
type LoginRequest struct {
	Account  string
	Password string
}

// AuthResponse 认证响应
type AuthResponse struct {
	UserID   uint
	Username string
	Email    string
	Nickname string
	Avatar   string
	Token    string
}

// UpdateUserRequest 更新用户信息请求
type UpdateUserRequest struct {
	Nickname string
	Avatar   string
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string
	NewPassword string
}
