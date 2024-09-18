package routes

import (
	"github.com/gin-gonic/gin"
	"iMessage/handlers"
	"iMessage/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(middleware.Logger()) // 使用自定义日志中间件
	r.POST("/upload", handlers.UploadFile)
	return r
}
