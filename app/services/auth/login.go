package auth

import (
	"jd_workout_golang/app/services/jwtHelper"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
	db "jd_workout_golang/lib/database"
	"jd_workout_golang/app/models"
)

func LoginAction(c *gin.Context) {
	loginForm := struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&loginForm); err != nil {
		c.JSON(422, gin.H{
			"message": "缺少必要欄位",
		})

		return
	}

	db := db.InitDatabase()

	user := models.User{}

	if _, err := validateLogin(&user, db); err != nil {
		c.JSON(422, gin.H{
			"message": "帳號或密碼錯誤",
			"error":   err.Error(),
		})

		return
	}
	
	token, _ := jwtHelper.GenerateToken(&user)

	c.JSON(200, gin.H{
		"message": "login success",
		"token":   token,
	})
}

func validateLogin(User *models.User, db *gorm.DB) (bool, error) {
	result := db.Where("email = ?", User.Email).First(&User)

	if result.Error != nil || result.RowsAffected == 0 {
		return false, result.Error
	}

	error := bcrypt.CompareHashAndPassword([]byte(User.Password), []byte(User.Password))

	if error != nil {
		return false, error
	}

	return true, nil
}
