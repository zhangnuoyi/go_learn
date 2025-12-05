package main

import (
	"fmt"
	"log"
	"net"

	"github.com/user/rpc/models"
	"github.com/user/rpc/proto"
	"github.com/user/rpc/services"
	"google.golang.org/grpc"
)

func main() {
	// 创建用户存储实例
	userStore := models.NewInMemoryUserStore()

	// 创建用户服务实例
	userService := services.NewUserService(userStore)

	// 创建 gRPC 服务器
	server := grpc.NewServer()

	// 注册用户服务
	proto.RegisterUserServiceServer(server, userService)

	// 监听端口
	port := 50051
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("无法监听端口 %d: %v", port, err)
	}

	log.Printf("gRPC 服务启动，监听端口 %d", port)

	// 启动服务器
	if err := server.Serve(listener); err != nil {
		log.Fatalf("gRPC 服务启动失败: %v", err)
	}
}
