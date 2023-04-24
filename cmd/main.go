package main

import (
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	env "github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"jd_workout_golang/app/middleware"
	docs "jd_workout_golang/docs"
	"jd_workout_golang/internal/router"
	"log"
	"os"
	"time"
)

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type Bearer followed by a space and JWT token.
// @scope.write Grants write access
func main() {
	env.Load()

	sentryDsn := os.Getenv("SENTRY_DSN")
	if sentryDsn != "" {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              sentryDsn,
			TracesSampleRate: 1.0,
		})

		if err != nil {
			log.Fatalf("sentry.Init: %s", err)
		}

		defer sentry.Flush(2 * time.Second)

		sentry.CaptureMessage("手動發送 Test sentry error - type 3")
	}

	r := SetupRouter()
	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(":80") // listen and serve on
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
