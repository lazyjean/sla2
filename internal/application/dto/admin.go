package dto

import (
	"github.com/lazyjean/sla2/internal/domain/entity"
)

// InitializeSystemRequest 初始化系统请求
type InitializeSystemRequest struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Nickname string `json:"nickname" validate:"required,max=50"`
	Password string `json:"password" validate:"required,min=8"`
	Email    string `json:"email" validate:"required,email"`
}

// InitializeSystemResponse 初始化系统响应
type InitializeSystemResponse struct {
	Admin        *entity.Admin `json:"admin"`
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
}

// AdminLoginRequest 管理员登录请求
type AdminLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// AdminLoginResponse 管理员登录响应
type AdminLoginResponse struct {
	Admin        *entity.Admin `json:"admin"`
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
}
