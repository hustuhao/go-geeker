package middleware

import (
	"github.com/gin-gonic/gin"
)

func CheckAdminLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否登录
		c.Next()
	}
}
