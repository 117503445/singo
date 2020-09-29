package api

import (
	"github.com/gin-gonic/gin"
)

// Ping 状态检查页面
func Ping(c *gin.Context) {
	c.JSON(200, "pong")
}
