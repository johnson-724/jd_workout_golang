package auth

import (
	db "jd_workout_golang/lib/database"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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
	
	error := bcrypt.CompareHashAndPassword([] byte(user.Password), []byte(loginForm.Password))

	if result.Error != nil || result.RowsAffected == 0|| error != nil {
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
