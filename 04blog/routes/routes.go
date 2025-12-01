package routes

import (
	"github.com/gin-gonic/gin"
)

// gin routes 设置
func SetRoutes(r *gin.Engine) {
	//后续接口路由
	//添加一个hello路由
	r.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})
}
