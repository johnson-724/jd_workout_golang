package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"jd_workout_golang/app/services/jwtHelper"
	"net/http"
)

var Uid uint

func ValidateToken(c *gin.Context) {
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
			"message": message,
			"error": message,
		})

		c.Abort()

		return
	}

	fmt.Println(Uid)

	c.Next()
}
