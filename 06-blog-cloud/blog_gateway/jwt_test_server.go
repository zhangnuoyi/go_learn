package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gateway/utils"
)

// 模拟用户信息
var mockUsers = map[string]string{
	"test": "123456", // 用户名: 密码
}

// 登录请求结构
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// 响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// 登录处理函数
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var loginReq LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"code": 400, "message": "请求参数错误"}`))
		return
	}

	// 验证用户名和密码
	if storedPassword, ok := mockUsers[loginReq.Username]; !ok || storedPassword != loginReq.Password {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"code": 401, "message": "用户名或密码错误"}`))
		return
	}

	// 生成JWT令牌
	token, err := utils.GenerateToken(1, loginReq.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"code": 500, "message": "生成令牌失败"}`))
		return
	}

	// 返回令牌
	response := Response{
		Code:    200,
		Message: "登录成功",
		Data: map[string]string{
			"token": token,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// 受保护的资源
func protectedHandler(w http.ResponseWriter, r *http.Request) {
	// 从请求头获取令牌
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"code": 401, "message": "缺少令牌"}`))
		return
	}

	// 解析令牌
	claims, err := utils.ParseToken(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"code": 401, "message": "无效的令牌"}`))
		return
	}

	// 返回受保护的数据
	response := Response{
		Code:    200,
		Message: "获取受保护资源成功",
		Data: map[string]interface{}{
			"userID":   claims.UserID,
			"username": claims.Username,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	// 注册路由
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/protected", protectedHandler)

	// 启动服务器
	fmt.Println("JWT测试服务器启动在 http://localhost:8089")
	log.Fatal(http.ListenAndServe(":8089", nil))
}
