package middleware

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"

	"gateway/utils"
)

// JwtMiddleware JWT鉴权中间件
type JwtMiddleware struct {
	// 不需要认证的路径列表
	ExcludePaths []string
}

// NewJwtMiddleware 创建JWT中间件实例
func NewJwtMiddleware(excludePaths []string) *JwtMiddleware {
	return &JwtMiddleware{
		ExcludePaths: excludePaths,
	}
}

// Handle 实现中间件处理函数
func (m *JwtMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 检查当前路径是否需要排除JWT验证
		path := r.URL.Path
		if m.shouldExclude(path) {
			next(w, r)
			return
		}

		// 从请求头获取token
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"code": 401, "message": "缺少Authorization请求头"}`))
			return
		}

		// 解析Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"code": 401, "message": "Authorization请求头格式错误"}`))
			return
		}

		tokenString := parts[1]

		// 验证token
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			log.Printf("Token验证失败: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"code": 401, "message": "无效的token"}`))
			return
		}

		// 将用户信息添加到请求上下文中
		ctx := r.Context()
		ctx = context.WithValue(ctx, "userID", claims.UserID)
		ctx = context.WithValue(ctx, "username", claims.Username)

		// 将用户信息添加到请求头中，以便传递给上游服务
		r = r.WithContext(ctx)
		r.Header.Set("X-User-ID", strconv.FormatInt(claims.UserID, 10))
		r.Header.Set("X-Username", claims.Username)

		// 继续处理请求
		next(w, r)
	}
}

// shouldExclude 判断路径是否需要排除JWT验证
func (m *JwtMiddleware) shouldExclude(path string) bool {
	for _, excludePath := range m.ExcludePaths {
		// 精确匹配或前缀匹配
		if path == excludePath || strings.HasPrefix(path, excludePath+"/") {
			return true
		}
	}
	return false
}

// MiddlewareFunc 返回中间件函数
func (m *JwtMiddleware) MiddlewareFunc() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			m.Handle(next.ServeHTTP)(w, r)
		})
	}
}
