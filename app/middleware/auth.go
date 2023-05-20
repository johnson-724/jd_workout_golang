package middleware

import (
	repo "jd_workout_golang/app/repositories/user"
	"jd_workout_golang/app/services/jwtHelper"
	"net/http"

	"github.com/gin-gonic/gin"
)

var Uid uint

func ValidateToken(c *gin.Context) {
	val := c.GetHeader("Authorization")

	if val == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "JWT token is empty",
		})

		c.Abort()

		return
	}

	res, ok := jwtHelper.ParseJwtToken(val, &Uid)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": res.Error.Error(),
			"error":   res.Error.Error(),
		})

		c.Abort()
		
		return
	}
	
	user, _ := repo.GetUserById(Uid)
	
	// password forgeted but not reset
	// fresh token
	if user.ResetPassword == 1 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "JWT token is invalid",
			"error":   "JWT token is invalid",
		})

		c.Abort()

		return
	}

	if res.ResetPassword == 1{
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "密碼已修改，請重新登入",
			"error":   "密碼已修改，請重新登入",
		})

		c.Abort()

		return
	}


	c.Next()
}

func VaildateResetPassword(c *gin.Context) {
	val := c.GetHeader("Authorization")

	if val == "" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "JWT token is empty",
		})

		c.Abort()

		return
	}

	res, ok := jwtHelper.ParseJwtToken(val, &Uid)

	if !ok {
		c.JSON(http.StatusForbidden, gin.H{
			"message": res.Error.Error(),
			"error":   res.Error.Error(),
		})

		c.Abort()

		return
	}

	c.Next()
}
