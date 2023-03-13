package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/gin-gonic/gin"
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

	user := user{
		Username: registerForm.Username,
		Email:    registerForm.Email,
		Password: hex.EncodeToString(sha256.New().Sum([]byte(registerForm.Password))),
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
