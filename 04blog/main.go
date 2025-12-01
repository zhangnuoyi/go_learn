package main

import (
	"04blog/config"
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
	if err := config.InitDB(); err != nil {
		fmt.Printf("数据库初始化失败: %v", err)
	}

	//创建gin路由
	r := gin.Default()
	//设置路由
	routes.SetRoutes(r)
	//启动服务器
	r.Run(fmt.Sprintf("%s:%s", config.AppConfig.Server.Host, config.AppConfig.Server.Port))

}
