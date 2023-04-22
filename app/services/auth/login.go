package auth

import (
	"jd_workout_golang/app/models"
	"jd_workout_golang/app/services/jwtHelper"
	db "jd_workout_golang/lib/database"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// LoginAction logs in a user with the provided email and password,
// and generates a JWT token for the user.
//
// @Summary Login user
// @Description Logs in a user with the provided email and password, and generates a JWT token for the user
// @Tags Auth
// @Accept x-www-form-urlencoded	
// @Produce json
// @Param email formData string true "User email"
// @Param password formData string true "User password"
// @Success 200 {string} string "{'message': 'login success', 'token': 'JWT token'}"
// @Failure 422 {string} string "{'message': '缺少必要欄位', 'error': 'error message'}"
// @Failure 422 {string} string "{'message': '帳號或密碼錯誤', 'error': 'error message'}"
// @Router /login [post]
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
