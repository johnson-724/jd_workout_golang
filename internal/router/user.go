package router

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func RegisterUser(r *gin.RouterGroup) {
	r.GET("/user", func(c *gin.Context) {
		c.String(200, "user")
	})

	r.POST("/user/register", registerAction)

	r.POST("/user/login", loginAction)
}

type registerForm struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type user struct {
	gorm.Model
	Username string
	Email    string
	Password string
}

func registerAction(c *gin.Context) {
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

	storeUser(&user)

	c.JSON(200, gin.H{
		"message": "register success",
	})
}

func loginAction(c *gin.Context) {
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

	db := initDatabase()

	user := user{}
	result := db.Where("email = ?", loginForm.Email).First(&user)
	password := hex.EncodeToString(sha256.New().Sum([]byte(loginForm.Password)))

	if result.Error != nil || result.RowsAffected == 0|| user.Password != password {
		c.JSON(422, gin.H{
			"message": "帳號或密碼錯誤",
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "login success",
	})
}

func storeUser(u *user) *user {
	db := initDatabase()

	db.Create(u)

	return u
}

func initDatabase() *gorm.DB {
	db, err := gorm.Open(
		mysql.New(mysql.Config{
			DSN: "root:root@tcp(mysql:3306)/jd_workout?charset=utf8mb4&parseTime=True&loc=Local",
		}),
		&gorm.Config{},
	)

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(user{})

	return db
}
