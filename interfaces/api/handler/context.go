package handler

import (
	"github.com/gin-gonic/gin"
)

// 用户ID的上下文键名
const UserIDKey = "user_id"

// getUserIDFromContext 从上下文中获取用户ID
func getUserIDFromContext(c *gin.Context) uint {
	// 从上下文中获取用户ID
	if userID, exists := c.Get(UserIDKey); exists {
		if id, ok := userID.(uint); ok {
			return id
		}
	}
	return 0 // 如果获取失败返回 0
}

// SetUserIDToContext 设置用户ID到上下文
func SetUserIDToContext(c *gin.Context, userID uint) {
	c.Set(UserIDKey, userID)
}
