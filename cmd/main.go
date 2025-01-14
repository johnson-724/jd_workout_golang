package main

import (
	"jd_workout_golang/app/middleware"
	authAction "jd_workout_golang/app/services/auth"
	docs "jd_workout_golang/docs"
	"jd_workout_golang/internal/router"
	"jd_workout_golang/lib/database"
	"jd_workout_golang/lib/file"
	"log"
	"os"
	"time"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type Bearer followed by a space and JWT token.
// @scope.write Grants write access
func main() {
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

func init() {
	file.LoadConfigAndEnv()
	database.InitDatabase()
}

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors())

	r.GET("/verify-email", authAction.VerifyEmail)

	// 註冊 router group
	apiGroup := r.Group("/api/v1")
	// 註冊 user router
	router.RegisterUser(apiGroup)

	router.RegisterEquip(apiGroup)

	router.RegisterRecord(apiGroup)

	return r
}
