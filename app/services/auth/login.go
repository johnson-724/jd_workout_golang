package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	db "jd_workout_golang/lib/database"
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

	user := user{}
	result := db.Where("email = ?", loginForm.Email).First(&user)
	password := hex.EncodeToString(sha256.New().Sum([]byte(loginForm.Password)))

	if result.Error != nil || result.RowsAffected == 0|| user.Password != password {
		c.JSON(422, gin.H{
			"message": "帳號或密碼錯誤",
			"error": result.Error,
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "login success",
	})
}
