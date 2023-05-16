package middleware

import (
	"github.com/gin-gonic/gin"
	"jd_workout_golang/lib/redis"
	"net/http"
)

func ApiRateLimit(c *gin.Context) {

	err := redis.ApiRateLimit(c) 

	if err != nil {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": "Too many requests",
		})

		c.Abort()

		return
	}

	c.Next()
}