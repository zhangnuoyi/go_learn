package api

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/user/rpc/models"
	"github.com/user/rpc/services"
)

// Server API 服务器
type Server struct {
	router     *gin.Engine
	userHandler *UserHandler
}

// NewServer 创建 API 服务器实例
func NewServer() *Server {
	// 创建用户存储实例
	userStore := models.NewInMemoryUserStore()

	// 创建用户服务实例
	userService := services.NewUserService(userStore)

	// 创建用户 API 处理器
	userHandler := NewUserHandler(userService)

	// 创建 Gin 路由器
	router := gin.Default()

	return &Server{
		router:     router,
		userHandler: userHandler,
	}
}

// SetupRoutes 设置路由
func (s *Server) SetupRoutes() {
	// API 根路径
	api := s.router.Group("/api")
	{
		// 用户相关路由
		users := api.Group("/users")
		{
			users.POST("", s.userHandler.CreateUser)
			users.GET("/:id", s.userHandler.GetUser)
			users.PUT("/:id", s.userHandler.UpdateUser)
			users.DELETE("/:id", s.userHandler.DeleteUser)
			users.GET("", s.userHandler.ListUsers)
		}
	}

	// 健康检查
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "服务运行正常",
		})
	})
}

// Start 启动服务器
func (s *Server) Start(port int) error {
	// 设置路由
	s.SetupRoutes()

	// 启动服务器
	addr := fmt.Sprintf(":%d", port)
	log.Printf("RESTful API 服务启动，监听端口 %d", port)
	return s.router.Run(addr)
}
