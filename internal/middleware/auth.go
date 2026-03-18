package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"nft-marketplace/pkg/jwt"
	"nft-marketplace/pkg/response"
)

func Auth(jwtSecret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Header 获取 Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, 401, "Authorization header required")
			c.Abort()
			return
		}

		// 提取 Token（Bearer <token>）
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, 401, "Invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 验证 Token
		claims, err := jwt.ParseToken(tokenString, jwtSecret)
		if err != nil {
			response.Error(c, 401, "Invalid token")
			c.Abort()
			return
		}

		// 将用户信息存储到 Context
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}
