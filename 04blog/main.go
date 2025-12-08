package main

import (
	"04blog/config"
	"04blog/middleware"
	"04blog/routes"
	"fmt"

	"github.com/gin-gonic/gin" // 需运行 go mod tidy 安装依赖
)

func main() {

	//加载服务器配置信息
	if err := config.LoadConfig(); err != nil {
		fmt.Printf("加载配置失败: %v", err)
	}

	//数据库初始化
	db, err := config.InitDB()
	if err != nil {
		fmt.Printf("数据库初始化失败: %v", err)
	}

	//Redis初始化
	redisClient, err := config.InitRedis()
	if err != nil {
		fmt.Printf("Redis初始化失败: %v", err)
	}

	//创建gin路由
	r := gin.Default()
	//添加全局错误处理中间件
	r.Use(middleware.ErrorHandler())

	// 添加自定义日志中间件
	r.Use(middleware.Logger())
	//添加认证中间件
	r.Use(middleware.AuthMiddleware())
	//设置路由
	routes.SetRoutes(r, db, redisClient)
	//启动服务器
	r.Run(fmt.Sprintf("%s:%s", config.AppConfig.Server.Host, config.AppConfig.Server.Port))

}
