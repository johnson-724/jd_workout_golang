package router

import (
	"github.com/gin-gonic/gin"
	"jd_workout_golang/app/services/auth"
)

func RegisterUser(r *gin.RouterGroup) {
	r.GET("/user", func(c *gin.Context) {
		c.String(200, "user")
	})

	r.POST("/user/register", auth.RegisterAction)

	r.POST("/user/login", auth.LoginAction)
}