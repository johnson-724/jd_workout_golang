package auth

import (
	"github.com/gin-gonic/gin"
	"jd_workout_golang/app/middleware"
)

func InfoAction(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "ok",
		"data": gin.H{
			"id": middleware.Uid,
		},
	})
}
