package dto

// RegisterDTO 注册请求DTO
type RegisterDTO struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Nickname string `json:"nickname"`
}

// LoginDTO 登录请求DTO
type LoginDTO struct {
	Account  string `json:"account" binding:"required"`  // 账号（用户名或邮箱）
	Password string `json:"password" binding:"required"` // 密码
}

// TokenDTO 令牌DTO
type TokenDTO struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}
