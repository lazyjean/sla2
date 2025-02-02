package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lazyjean/sla2/interfaces/api/handler"
	"github.com/lazyjean/sla2/pkg/jwt"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从不同来源获取 token
		var tokenString string

		// 1. 先尝试从 Authorization header 获取
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) > 7 && strings.ToUpper(authHeader[0:7]) == "BEARER " {
			tokenString = authHeader[7:]
		}

		// 2. 尝试从 query 参数获取
		if tokenString == "" {
			tokenString = c.Query("token")
		}

		// 4. 尝试从 cookie 获取
		if tokenString == "" {
			if cookie, err := c.Cookie("jwt"); err == nil {
				tokenString = cookie
			}
		}

		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		claims, err := jwt.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// 将用户ID存入上下文
		handler.SetUserIDToContext(c, uint(claims.UserID))
		c.Next()
	}
}
