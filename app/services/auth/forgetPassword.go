package auth

import (
	"github.com/gin-gonic/gin"
	"jd_workout_golang/app/models"
	repo "jd_workout_golang/app/repositories/user"
	email "jd_workout_golang/lib/Email"
	"jd_workout_golang/app/services/jwtHelper"
	"os"
	"fmt"
)

type ForgetPassword struct {
	Email string `json:"email" form:"email" binding:"required"`
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
	user.Password = ""

	sendForgetPassword(user)

	repo.Update(user)

	c.JSON(200, gin.H{
		"message": "forget password success",
	})
}

func sendForgetPassword(u *models.User) error {
	payload := map[string]interface{}{
		"method": "forgetPassword",
	}

	token,err := jwtHelper.GenerateTokenWithPayload(u, payload)

	if err != nil {
		return err
	}

	baseUrl := os.Getenv("APP_URL") + "/reset-password?email=%s&token=%s"

	mail := email.Email{
		FromName:  os.Getenv("EMAIL_FROM_NAME"),
		FromEmail: os.Getenv("EMAIL_FROM_ADDRESS"),
		ToEmail:   u.Email,
		ToName:    u.Username,
		Subject:   "JD Workout 密碼重置信",
		Content:   fmt.Sprintf("請點擊以下連結重置密碼: %s", fmt.Sprintf(baseUrl, u.Email, token)),
	}

	err = email.Send(mail)

	if err != nil {
		return err
	}

	return nil
}
