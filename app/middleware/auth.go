package middleware

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	env "github.com/joho/godotenv"
	"jd_workout_golang/app/services/jwtHelper"
)

var Uid uint

func ValidateToken(c *gin.Context) {
	env.Load()

	val := c.GetHeader("Authorization")

	if val == "" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "JWT token is empty",
		})

		c.Abort()

		return
	}

	message, ok := jwtHelper.ValidateToken(val, &Uid)

	if !ok {
		c.JSON(http.StatusForbidden, gin.H{
			"error": message,
		})

		c.Abort()

		return
	}

	fmt.Println(Uid)

	c.Next()
}
