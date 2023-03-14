package auth

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	db "jd_workout_golang/lib/database"
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

	user := user{
		Username: registerForm.Username,
		Email:    registerForm.Email,
		Password: string(hash),
	}

	db := db.InitDatabase()

	db.AutoMigrate(&user)

	storeUser(&user, db)

	c.JSON(200, gin.H{
		"message": "register success",
	})
}

func storeUser(u *user, db *gorm.DB) *user {
	db.Create(u)

	return u
}
