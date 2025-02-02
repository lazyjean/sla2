package dto

// RegisterDTO 注册请求DTO
type RegisterDTO struct {
	Username string `json:"username" binding:"required,min=3,max=50" example:"johndoe"`
	Password string `json:"password" binding:"required,min=6,max=50" example:"password123"`
	Email    string `json:"email" binding:"omitempty,email" example:"john@example.com"`
	Phone    string `json:"phone" binding:"omitempty,len=11" example:"13800138000"`
}

// LoginDTO 登录请求DTO
type LoginDTO struct {
	// 支持用户名/邮箱/手机号登录
	Account  string `json:"account" binding:"required" example:"johndoe"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// TokenDTO 令牌响应DTO
type TokenDTO struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIs..."`
	TokenType   string `json:"token_type" example:"Bearer"`
	ExpiresIn   int    `json:"expires_in" example:"3600"` // 过期时间(秒)
}

// UserDTO 用户信息DTO
type UserDTO struct {
	ID       uint   `json:"id" example:"1"`
	Username string `json:"username" example:"johndoe"`
	Email    string `json:"email,omitempty" example:"john@example.com"`
	Phone    string `json:"phone,omitempty" example:"13800138000"`
}
