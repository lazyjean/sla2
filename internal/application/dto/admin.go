package dto

import (
	"time"

	"github.com/lazyjean/sla2/internal/domain/entity"
)

// SystemStatusResponse 系统状态响应
type SystemStatusResponse struct {
	Initialized bool `json:"initialized"`
}

// InitializeSystemRequest 初始化系统请求
type InitializeSystemRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// InitializeSystemResponse 初始化系统响应
type InitializeSystemResponse struct {
	Admin        *AdminInfo `json:"admin"`
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
}

// AdminLoginRequest 管理员登录请求
type AdminLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AdminLoginResponse 管理员登录响应
type AdminLoginResponse struct {
	Admin        *AdminInfo `json:"admin"`
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
}

// AdminRefreshTokenResponse 管理员刷新令牌响应
type AdminRefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// AdminInfoResponse 管理员信息响应
type AdminInfoResponse struct {
	ID       entity.UID `json:"id"`
	Username string     `json:"username"`
	Nickname string     `json:"nickname"`
	Roles    []string   `json:"permissions"`
}

// AdminInfo 管理员信息
type AdminInfo struct {
	ID        entity.UID `json:"id"`
	Username  string     `json:"username"`
	Nickname  string     `json:"nickname"`
	Roles     []string   `json:"roles"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
