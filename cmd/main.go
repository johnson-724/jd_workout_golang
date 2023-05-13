package main

import (
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"jd_workout_golang/app/middleware"
	authAction "jd_workout_golang/app/services/auth"
	docs "jd_workout_golang/docs"
	"jd_workout_golang/internal/router"
	"jd_workout_golang/lib/database"
	"jd_workout_golang/lib/file"
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
	setFileLog()
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

		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered from panic: %v\n", r)
				sentry.CurrentHub().Recover(r)
				sentry.Flush(time.Second * 5)
			}
		}()
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

	apiGroup.GET("/app/version", getAppVersion)

	// 註冊 user router
	router.RegisterUser(apiGroup)

	router.RegisterEquip(apiGroup)

	router.RegisterRecord(apiGroup)

	return r
}

func setFileLog() {
	logFile, err := os.OpenFile("logfile.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	log.SetOutput(logFile)
}

// @Summary Get app version
// @Description Get app version
// @Tags App
// @Accept  json
// @Produce  json
// @Success 200 {object} string "{'latestVersion: 1.0.0', 'requiredVersion': '1.0.0'}"
// @Router /app/version [get]
func getAppVersion(c *gin.Context) {

	c.JSON(200, gin.H{
		"latestVersion":   os.Getenv("APP_VERSION"),
		"requiredVersion": os.Getenv("APP_REQUIRED_VERSION"),
	})
}
