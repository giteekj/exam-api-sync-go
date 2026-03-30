package middleware

import (
	"net/http"
	"strings"

	"exam-api-sync-go/common/jwt"
	"exam-api-sync-go/common/response"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			response.Fail(c, http.StatusUnauthorized, "未提供认证信息")
			c.Abort()
			return
		}

		parts := strings.SplitN(auth, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Fail(c, http.StatusUnauthorized, "认证信息格式错误")
			c.Abort()
			return
		}

		// 解析token
		payload, err := jwt.ParseToken(parts[1])
		if err != nil {
			response.Fail(c, http.StatusUnauthorized, "无效的token")
			c.Abort()
			return
		}

		// 将用户信息存储到上下文
		c.Set("user_id", payload.UserID)
		c.Set("username", payload.Username)
		c.Set("role", payload.Role)

		c.Next()
	}
}
