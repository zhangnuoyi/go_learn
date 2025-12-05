package main

import (
	"04blog/config"
	"04blog/middleware"
	"04blog/routes"
	"fmt"

	"github.com/gin-gonic/gin" // 需运行 go mod tidy 安装依赖
)

func main() {
	fmt.Println("程序开始执行")

	//加载服务器配置信息
	fmt.Println("正在加载配置...")
	if err := config.LoadConfig(); err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		return
	}
	fmt.Println("配置加载成功")

	//数据库初始化
	fmt.Println("正在初始化数据库...")
	db, err := config.InitDB()
	if err != nil {
		fmt.Printf("数据库初始化失败: %v\n", err)
		return
	}
	fmt.Println("数据库初始化成功")

	//Redis初始化
	fmt.Println("正在初始化Redis...")
	redisClient, err := config.InitRedis()
	if err != nil {
		fmt.Printf("Redis初始化失败: %v\n", err)
		return
	}
	fmt.Println("Redis初始化成功")

	//创建gin路由
	fmt.Println("正在创建Gin路由...")
	r := gin.Default()
	
	//添加全局错误处理中间件
	fmt.Println("正在添加中间件...")
	r.Use(middleware.ErrorHandler())

	// 添加自定义日志中间件
	r.Use(middleware.Logger())
	//添加认证中间件
	r.Use(middleware.AuthMiddleware())
	
	//设置路由
	fmt.Println("正在设置路由...")
	routes.SetRoutes(r, db, redisClient)
	
	//启动服务器
	addr := fmt.Sprintf("%s:%s", config.AppConfig.Server.Host, config.AppConfig.Server.Port)
	fmt.Printf("服务器正在启动，监听地址: %s\n", addr)
	r.Run(addr)

}
