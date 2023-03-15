package router

import (
	"github.com/gin-gonic/gin"
	"jd_workout_golang/app/services/auth"
)

func RegisterUser(r *gin.RouterGroup) {
	r.POST("/register", auth.RegisterAction)
	r.POST("/login", auth.LoginAction)

	userGroup := r.Group("/user").Use(auth.ValidateToken)

	userGroup.GET("/", auth.InfoAction)
}