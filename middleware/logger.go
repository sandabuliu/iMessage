package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		latency := endTime.Sub(startTime)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		uri := c.Request.RequestURI

		logMessage := fmt.Sprintf("| %3d | %13v | %15s | %s | %s |",
			statusCode, latency, clientIP, method, uri)
		fmt.Println(logMessage)
	}
}
