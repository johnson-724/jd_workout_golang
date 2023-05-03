package auth

import (
	"fmt"
	"jd_workout_golang/app/models"
	"jd_workout_golang/app/services/jwtHelper"
	"jd_workout_golang/lib/database"
	"github.com/gin-gonic/gin"
)

func VerifyEmail(c *gin.Context) {
	user := models.User{}
	uid := uint(0)
	
	err, res := jwtHelper.ValidateToken(fmt.Sprintf("Bearer %s", c.Query("token")), &uid)

	if !res {
		c.JSON(422, gin.H{
			"message": "驗證連結無效 token",
			"error": err,
		})

		c.Abort()

		return
	}

	result := database.Connection.Model(&user).
		Where("email = ?", c.Query("email")).
		Where("email_verified = ?", 0).
		Find(&user)

	if result.RowsAffected == 0 || user.ID != uid {
		c.JSON(422, gin.H{
			"message": "驗證連結無效 email",
		})

		c.Abort()

		return
	}

	database.Connection.Model(&user).Where("id = ?", uid).Update("email_verified", 1)

	c.JSON(200, gin.H{
		"message": "信箱驗證成功",
	})
}
