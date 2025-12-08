package middleware

import (
	"log"
	"net/http"
	"time"
)

// LogMiddleware 日志中间件
type LogMiddleware struct {}

// NewLogMiddleware 创建日志中间件实例
func NewLogMiddleware() *LogMiddleware {
	return &LogMiddleware{}
}

// MiddlewareFunc 返回中间件函数
func (m *LogMiddleware) MiddlewareFunc() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 记录请求开始时间
			startTime := time.Now()

			// 获取客户端IP
			clientIP := r.RemoteAddr

			// 创建响应包装器，用于记录响应状态码
			wrapper := &responseWrapper{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			// 处理请求
			next.ServeHTTP(wrapper, r)

			// 计算响应时间
			duration := time.Since(startTime)

			// 记录日志
			log.Printf("[GATEWAY] %s %s %d %s %s",
				r.Method,
				r.URL.Path,
				wrapper.statusCode,
				duration,
				clientIP,
			)
		})
	}
}

// responseWrapper 响应包装器，用于记录响应状态码
type responseWrapper struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader 重写WriteHeader方法，记录状态码
func (w *responseWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
