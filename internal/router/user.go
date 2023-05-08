package router

import (
	"github.com/gin-gonic/gin"
	auth "jd_workout_golang/app/middleware"
	authAction "jd_workout_golang/app/services/auth"
)

func RegisterUser(r *gin.RouterGroup) {
	r.POST("/register", authAction.RegisterAction)
	r.POST("/login", authAction.LoginAction)
	r.GET("/login/google", authAction.LoginWithGoogleAction)
	r.GET("/login/google/redirect", authAction.LoginWithGoogleAuthkAction)

	userGroup := r.Group("/user").Use(auth.ValidateToken)

	userGroup.GET("/", authAction.InfoAction)
}