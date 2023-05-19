package auth

import (
	"fmt"
	"jd_workout_golang/app/middleware"
	"jd_workout_golang/app/models"
	repo "jd_workout_golang/app/repositories/user"
	email "jd_workout_golang/lib/Email"
	"os"
	helper "jd_workout_golang/lib/helper"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type ForgetPassword struct {
	Email string `json:"email" form:"email" binding:"required"`
}

type ResetPassword struct {
	Passowrd        string `json:"password" form:"password" binding:"required"`
	NewPassword     string `json:"newPassword" form:"newPassword" binding:"required"`
	ConfirmPassword string `json:"confirmPassword" form:"confirmPassword" binding:"required"`
}

// ForgetPasswordAction forget password API
// @Summary forget password
// @Description user forget password
// @Tags Auth
// @Accept json
// @Produce json
// @Param forgetPassword body ForgetPassword true "forgetPassword"
// @Success 200 {string} json "{"message": "forget password success"}"
// @Failure 422 {string} json "{"message": "forget password failed", "error": "forget password form validation failed"}"
// @Failure 422 {string} json "{"message": "查無此 email"}"
// @BasePath
// @Router /forget-password [post]
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

	user.ResetPassword = 1
	user.Password = helper.RandString(8)

	sendForgetPassword(user)

	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hash)

	repo.Update(user)

	c.JSON(200, gin.H{
		"message": "forget password success",
	})
}

func ResetPassowrd(c *gin.Context) {
	restPassword := ResetPassword{}
	if err := c.ShouldBind(&restPassword); err != nil {
		c.JSON(422, gin.H{
			"message": "缺少欄位",
			"error": err.Error(),
		})

		c.Abort()

		return
	}

	if restPassword.NewPassword != restPassword.ConfirmPassword {
		c.JSON(422, gin.H{
			"message": "密碼不一致",
		})

		c.Abort()

		return
	}

	user, e := repo.GetUserById(middleware.Uid)

	if e != nil {
		c.JSON(422, gin.H{
			"message": "重置無效",
		})

		c.Abort()

		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(restPassword.NewPassword), bcrypt.DefaultCost)

	user.Password = string(hash)
	user.ResetPassword = 0

	repo.Update(user)

	c.JSON(200, gin.H{
		"message": "密碼修改成功",
	})
}

func sendForgetPassword(u *models.User) error {
	mail := email.Email{
		FromName:  os.Getenv("EMAIL_FROM_NAME"),
		FromEmail: os.Getenv("EMAIL_FROM_ADDRESS"),
		ToEmail:   u.Email,
		ToName:    u.Username,
		Subject:   "JD Workout 密碼重置信",
		Content:   fmt.Sprintf("重置密碼: %s\n請儘速登入並修改密碼", fmt.Sprintf(u.Password)),
	}

	err := email.Send(mail)

	if err != nil {
		return err
	}

	return nil
}
