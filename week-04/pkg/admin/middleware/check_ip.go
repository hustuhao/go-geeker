package middleware

import "github.com/gin-gonic/gin"

func CheckIP() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查访问ip是否在ip白名单中
		c.Next()
	}
}
