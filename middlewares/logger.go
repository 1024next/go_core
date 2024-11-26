package middlewares

import (
	"log"

	"github.com/gin-gonic/gin"
)

// Logger 中间件记录每个请求的日志
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 处理请求前的日志
		log.Printf("Started %s %s", c.Request.Method, c.Request.URL.Path)

		// 处理请求
		c.Next()

		// 处理请求后的日志
		log.Printf("Completed %s %s with status %d", c.Request.Method, c.Request.URL.Path, c.Writer.Status())
	}
}
