package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	// TODO: 添加认证服务依赖
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

// @Router       /v1/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	// TODO: 实现登录逻辑
	c.JSON(http.StatusOK, gin.H{
		"message": "登录功能待实现",
	})
}

// @Router       /v1/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	// TODO: 实现注册逻辑
	c.JSON(http.StatusOK, gin.H{
		"message": "注册功能待实现",
	})
}
