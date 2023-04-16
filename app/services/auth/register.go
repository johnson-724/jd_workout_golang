package auth

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	db "jd_workout_golang/lib/database"
	"jd_workout_golang/app/models"
)

type registerForm struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func RegisterAction(c *gin.Context) {
	registerForm := registerForm{}
	if err := c.ShouldBindJSON(&registerForm); err != nil {
		c.JSON(422, gin.H{
			"message": "register failed",
		})

		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(registerForm.Password), bcrypt.DefaultCost)

	user := &models.User{
		Username: registerForm.Username,
		Email:    registerForm.Email,
		Password: string(hash),
	}

	if state := validateRegister(user, db.InitDatabase()); !state {
		c.JSON(422, gin.H{
			"message": "Email 重複",
			"error": "duplicate email",
		})

		return
	}

	db := db.InitDatabase()

	storeUser(user, db)

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
