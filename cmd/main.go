package main

import (
	"github.com/gin-gonic/gin"
	"jd_workout_golang/internal/router"
)

func main() {
	r := SetupRouter()

	r.Run(":80") // listen and serve on
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 註冊 router group
	apiGroup := r.Group("/api")
	// 註冊 user router
	router.RegisterUser(apiGroup)

	return r
}
