package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func main() {
	r := gin.Default()
	r.GET("/index", Index)

	r.GET("/hello", Hello)
	r.Run(":8080")
}
func Index(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code: http.StatusOK,
		Msg:  "success",
		Data: "hello world",
	})
}
func Hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success",
		"data": gin.H{
			"name": "gin",
		},
	})
}
