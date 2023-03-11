package router

import (
	"github.com/gin-gonic/gin"
)

func RegisterUser(r *gin.RouterGroup) {
	r.GET("/user", func(c *gin.Context) {
		c.String(200, "user")
	})
}
