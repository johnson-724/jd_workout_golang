package auth

import (
	"fmt"
	"jd_workout_golang/app/models"
	"jd_workout_golang/app/services/jwtHelper"
	email "jd_workout_golang/lib/Email"
	db "jd_workout_golang/lib/database"
	"os"	
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type registerForm struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
	Email    string `json:"email" form:"email" binding:"required"`
}

// RegisterAction register API
// @Summary register
// @Description user register
// @Tags Auth
// @Accept json
// @Produce json
// @Param registerForm body registerForm true "registerForm"
// @Success 200 {string} json "{"message": "register success"}"
// @Failure 422 {string} json "{"message": "register failed", "error": "register form validation failed"}"
// @Failure 422 {string} json "{"message": "Email 重複", "error": "duplicate email"}"
// @Failure 422 {string} json "{"message": "信箱尚未驗證"}"
// @Failure 500 {string} json "{"message": "register failed", "error": "server error"}"
// @Router /register [post]
func RegisterAction(c *gin.Context) {
	registerForm := registerForm{}
	if err := c.ShouldBind(&registerForm); err != nil {
		c.JSON(422, gin.H{
			"message": "register failed",
			"error":   err.Error(),
		})

		c.Abort()

		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(registerForm.Password), bcrypt.DefaultCost)

	user := &models.User{
		Username: registerForm.Username,
		Email:    registerForm.Email,
		Password: string(hash),
	}

	if state := validateRegister(user, db.Connection); !state {
		c.JSON(422, gin.H{
			"message": "Email 重複",
			"error":   "duplicate email",
		})

		c.Abort()

		return
	}

	storeUser(user, db.Connection)

	err := generateVerifyEmail(user)

	if err != nil {
		c.JSON(500, gin.H{
			"message": "register failed",
			"error":   "server error",
		})

		c.Abort()

		return
	}

	c.JSON(200, gin.H{
		"message": "register success",
	})
}

func validateRegister(u *models.User, db *gorm.DB) bool {
	result := db.Where("email = ?", u.Email).First(&u)

	// duplicate email
	if result.RowsAffected == 1 {
		return false
	}

	return true
}

func storeUser(u *models.User, db *gorm.DB) *models.User {
	db.Create(u)

	return u
}

func generateVerifyEmail(user *models.User) *error {
	token, _ := jwtHelper.GenerateToken(user)
	baseUrl := os.Getenv("APP_URL") + "/verifyEmail?token=%s"

	mail := email.Email{
		FromName:  os.Getenv("EMAIL_FROM_NAME"),
		FromEmail: os.Getenv("EMAIL_FROM_ADDRESS"),
		ToEmail:   user.Email,
		ToName:    user.Username,
		Subject:   "JD Workout 驗證信",
		Content:   fmt.Sprintf("請點擊以下連結驗證信箱: %s", fmt.Sprintf(baseUrl, token)),
	}

	err := email.Send(mail)

	if err != nil {
		return err
	}

	return nil
}
