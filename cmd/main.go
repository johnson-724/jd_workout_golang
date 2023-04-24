package main

import (
	"os"
	"jd_workout_golang/lib/file"
	"jd_workout_golang/lib/database"
	"jd_workout_golang/internal/router"
	"jd_workout_golang/app/middleware"
	"github.com/gin-gonic/gin"
	docs "jd_workout_golang/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type Bearer followed by a space and JWT token.
// @scope.write Grants write access
func main() {
	r := SetupRouter()
	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(":80") // listen and serve on
}

func init() {
	file.LoadConfigAndEnv()
	database.InitDatabase()
}

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors())

	// 註冊 router group
	apiGroup := r.Group("/api/v1")
	// 註冊 user router
	router.RegisterUser(apiGroup)

	router.RegisterEquip(apiGroup)

	return r
}
