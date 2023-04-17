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
		Email    string `json:"email" form:"email" binding:"required"`
		Password string `json:"password" form:"password" binding:"required"`
	}{}

	if err := c.ShouldBind(&loginForm); err != nil {
		c.JSON(422, gin.H{
			"message": "缺少必要欄位",
			"error":   err.Error(),
		})

		return
	}

	db := db.InitDatabase()

	user := models.User{
		Email:    loginForm.Email,
		Password: loginForm.Password,
	}

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

func validateLogin(user *models.User, db *gorm.DB) (bool, error) {
	record := models.User{}

	result := db.Where("email = ?", user.Email).First(&record)

	if result.Error != nil || result.RowsAffected == 0 {
		return false, result.Error
	}

	error := bcrypt.CompareHashAndPassword([]byte(record.Password), []byte(user.Password))

	if error != nil {
		return false, error
	}

	*user = record

	return true, nil
}
