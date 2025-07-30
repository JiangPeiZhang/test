package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthResponse 健康检查响应结构
type HealthResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func main() {
	// 创建Gin引擎
	r := gin.Default()

	// 注册Health接口 - POST请求
	r.POST("/health", func(c *gin.Context) {
		response := HealthResponse{
			Code: 0,
			Msg:  "",
		}
		c.JSON(http.StatusOK, response)
	})

	// 启动服务器，监听8080端口
	r.Run(":8080")
}
