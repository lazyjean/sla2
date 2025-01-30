package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lazyjean/sla2/config"
	"github.com/lazyjean/sla2/interfaces/api/handler"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 这里实现您的认证逻辑
		// 例如解析 JWT token 获取用户ID

		userID := extractUserIDFromToken(c) // 实现此函数以从token中提取用户ID

		// 设置用户ID到上下文
		handler.SetUserIDToContext(c, userID)

		c.Next()
	}
}

func extractUserIDFromToken(c *gin.Context) uint {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		return 0
	}

	// 移除 "Bearer " 前缀
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// 解析 JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		// 从配置中读取密钥
		return []byte(config.GetConfig().JWT.SecretKey), nil
	})

	if err != nil || !token.Valid {
		return 0
	}

	// 从 token 中提取用户 ID
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if userID, exists := claims["userId"]; exists {
			// 将 interface{} 转换为 uint
			if id, ok := userID.(float64); ok {
				return uint(id)
			}
		}
	}

	return 0
}
