package auth

import (
	"github.com/gin-gonic/gin"
	repo "jd_workout_golang/app/repositories/user"
)

type ForgetPassword struct {
	Email string `json:"email" form:"email" binding:"required"`
}

func ForgetPasswordAction(c *gin.Context) {
	forgetPassword := ForgetPassword{}

	if err := c.ShouldBind(&forgetPassword); err != nil {
		c.JSON(422, gin.H{
			"message": "缺少 email 欄位",
			"error":   err.Error(),
		})

		c.Abort()

		return
	}

	user, err := repo.GetUserByEmail(forgetPassword.Email)

	if err != nil {
		c.JSON(422, gin.H{
			"message": "查無此 email",
			"error":   err.Error(),
		})

		c.Abort()

		return
	}

	user.ResetPassword = 0
	user.Password = ""

	sendForgetPassword()

	repo.Update(user)
}

func sendForgetPassword() {

}
