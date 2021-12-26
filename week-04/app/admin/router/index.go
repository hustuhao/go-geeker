package router

import "github.com/gin-gonic/gin"

const INDEX = `请登录`

func AdminIndex(c *gin.Context) {
	c.String(200, INDEX)
}
