package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func response(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(code, Response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

var mapCode = map[int]string{
	1001: "权限错误",
	1002: "角色错误",
}

func Ok(c *gin.Context, data interface{}) {
	response(c, http.StatusOK, "success", data)
}
func OkMsg(c *gin.Context, msg string) {
	response(c, http.StatusOK, msg, map[string]interface{}{})
}
func OkData(c *gin.Context, data interface{}) {
	response(c, http.StatusOK, "success", data)
}

// 失败
func Fail(c *gin.Context) {
	response(c, 500, "fail", map[string]interface{}{})
}
func FailMsg(c *gin.Context, code int, msg string) {

	mapMsg, ok := mapCode[code]
	if !ok {
		mapMsg = "服务错误，请稍后再试"
	}
	response(c, code, mapMsg, map[string]interface{}{})
}
