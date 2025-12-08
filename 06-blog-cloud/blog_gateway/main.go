package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gateway/config"
	"gateway/middleware"
	"gateway/router"
	"gateway/utils"
)

// EnhancedLoginResponse 增强的登录响应结构
type EnhancedLoginResponse struct {
	Token string `json:"token"`
}

func main() {
	// 1. 加载配置文件
	cfg, err := config.LoadConfig("etcd/gateway.yaml")
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	// 2. 初始化路由管理器
	router := router.NewRouter(cfg)

	// 3. 初始化中间件
	jwtMiddleware := middleware.NewJwtMiddleware(cfg.Jwt.ExcludePaths)
	logMiddleware := middleware.NewLogMiddleware()

	// 4. 自定义处理器
	handler := http.NewServeMux()

	// 注册登录路由处理函数（保持原有逻辑）
	handler.HandleFunc("/v1/user/login", func(w http.ResponseWriter, r *http.Request) {
		// 只处理POST请求
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(`{"code": 405, "message": "Method not allowed"}`))
			return
		}

		// 创建上下文，设置用户ID（key为user_id）
		ctx := context.WithValue(r.Context(), "user_id", int64(7))
		r = r.WithContext(ctx)

		// 直接返回增强的登录响应
		w.Header().Set("Content-Type", "application/json")

		// 生成JWT token
		token, err := utils.GenerateToken(7, "zq")
		if err != nil {
			token = "default_token_for_testing"
		}

		response := EnhancedLoginResponse{
			Token: token,
		}

		json.NewEncoder(w).Encode(response)
	})

	// 处理所有其他请求，根据配置文件进行路由匹配和转发
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 匹配路由
		upstreamName := router.MatchRoute(r.Method, r.URL.Path)
		if upstreamName == "" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"code": 404, "message": "路由不存在"}`))
			return
		}

		// 获取对应的代理
		proxy := router.GetProxy(upstreamName)
		if proxy == nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"code": 500, "message": "代理服务不可用"}`))
			return
		}

		// 转发请求
		proxy.ServeHTTP(w, r)
	})

	// 5. 构建中间件链
	var finalHandler http.Handler = handler
	finalHandler = logMiddleware.MiddlewareFunc()(finalHandler)
	finalHandler = jwtMiddleware.MiddlewareFunc()(finalHandler)

	// 6. 启动HTTP服务器
	serverAddr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	server := &http.Server{
		Addr:    serverAddr,
		Handler: finalHandler,
	}

	log.Printf("网关服务启动中，监听地址: %s", serverAddr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("启动网关失败: %v", err)
	}
}
