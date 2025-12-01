package middleware

import (
	"04blog/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查请求头是否包含 Authorization 字段
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
			// 未授权，终止处理
			c.Abort()
			return
		}
		// 验证 Authorization 值是否正确 查看token是否以Bearer开头
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(401, gin.H{"error": "无效的授权令牌"})
			c.Abort()
			return
		}

		//解析token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的授权令牌"})
			c.Abort()
			return
		}
		//设置用户ID到上下文
		c.Set("user_id", claims.UserID)
		//设置用户名到上下文
		c.Set("username", claims.Username)

		// 授权成功，继续处理请求
		c.Next()
	}
}
