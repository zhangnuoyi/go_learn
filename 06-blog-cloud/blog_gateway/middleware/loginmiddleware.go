package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"gateway/utils"
)

// GenerateTokenForLogin 生成登录用的JWT令牌
func GenerateTokenForLogin(userID int64, username string) (string, error) {
	return utils.GenerateToken(userID, username)
}

// LoginResponse 用户登录响应结构
type LoginResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// EnhancedLoginResponse 增强的登录响应结构
type EnhancedLoginResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Token string `json:"token"`
}

// LoginMiddleware 登录处理器中间件
type LoginMiddleware struct {
	// 代理请求到的目标URL
	TargetURL string
}

// NewLoginMiddleware 创建登录处理器中间件实例
func NewLoginMiddleware(targetURL string) *LoginMiddleware {
	return &LoginMiddleware{
		TargetURL: targetURL,
	}
}

// Handle 实现中间件处理函数
func (m *LoginMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 只处理登录请求
		if r.URL.Path != "/v1/user/login" || r.Method != http.MethodPost {
			next(w, r)
			return
		}

		// 读取原始请求体
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"code": 500, "message": "读取请求体失败"}`))
			return
		}

		// 恢复请求体，以便后续使用
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// 创建一个自定义的ResponseWriter来捕获上游服务的响应
		recorder := &responseRecorder{
			header: make(http.Header),
		}

		// 调用下一个处理器（通常是代理到用户API服务）
		next(recorder, r)

		// 如果上游服务返回的不是200 OK，直接返回原始响应
		if recorder.statusCode != http.StatusOK {
			w.WriteHeader(recorder.statusCode)
			w.Write(recorder.body)
			return
		}

		// 解析上游服务返回的用户信息
		var loginResp LoginResponse
		if err := json.Unmarshal(recorder.body, &loginResp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"code": 500, "message": "解析用户信息失败"}`))
			return
		}

		// 生成JWT Token
		token, err := utils.GenerateToken(loginResp.ID, loginResp.Name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"code": 500, "message": "生成Token失败"}`))
			return
		}

		// 创建增强的响应，包含Token
		enhancedResp := EnhancedLoginResponse{
			ID:    loginResp.ID,
			Name:  loginResp.Name,
			Email: loginResp.Email,
			Token: token,
		}

		// 将用户ID添加到上下文中（如果后续中间件需要使用）
		ctx := r.Context()
		ctx = context.WithValue(ctx, "user_id", loginResp.ID)
		r = r.WithContext(ctx)

		// 将用户ID添加到请求头中
		r.Header.Set("X-User-ID", strconv.FormatInt(loginResp.ID, 10))

		// 返回增强的响应
		w.Header().Set("Content-Type", "application/json")
		respData, _ := json.Marshal(enhancedResp)
		w.Write(respData)
	}
}

// MiddlewareFunc 返回中间件函数
func (m *LoginMiddleware) MiddlewareFunc() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			m.Handle(next.ServeHTTP)(w, r)
		})
	}
}

// responseRecorder 用于捕获HTTP响应的自定义ResponseWriter
type responseRecorder struct {
	header     http.Header
	body       []byte
	statusCode int
}

// Header 实现http.ResponseWriter接口
func (r *responseRecorder) Header() http.Header {
	return r.header
}

// Write 实现http.ResponseWriter接口
func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body = append(r.body, b...)
	return len(b), nil
}

// WriteHeader 实现http.ResponseWriter接口
func (r *responseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
}
