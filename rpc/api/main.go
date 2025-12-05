package main

import (
	"log"

	"github.com/user/rpc/api"
)

func main() {
	// 创建 API 服务器实例
	server := api.NewServer()

	// 启动服务器
	port := 8080
	log.Printf("准备启动 RESTful API 服务，端口: %d", port)
	if err := server.Start(port); err != nil {
		log.Fatalf("RESTful API 服务启动失败: %v", err)
	}
}
