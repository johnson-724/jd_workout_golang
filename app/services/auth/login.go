package auth

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	env "github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	db "jd_workout_golang/lib/database"
	"os"
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

	error := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginForm.Password))

	if result.Error != nil || result.RowsAffected == 0 || error != nil {
		c.JSON(422, gin.H{
			"message": "帳號或密碼錯誤",
			"error":   result.Error,
		})

		return
	}
	
	token, _ := generateToken(&user)

	c.JSON(200, gin.H{
		"message": "login success",
		"token":   token,
	})
}

func generateToken(u *user) (string, error) {
	env.Load()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": u.ID,
	})

	return token.SignedString([]byte(os.Getenv("APP_KEY")))
}
