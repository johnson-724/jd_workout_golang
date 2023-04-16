package jwtHelper

import (
	jwt "github.com/golang-jwt/jwt/v5"
	env "github.com/joho/godotenv"
	"os"
	"jd_workout_golang/app/models"
)

func GenerateToken (u *models.User) (string, error) {
	env.Load()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": u.ID,
	})

	return token.SignedString([]byte(os.Getenv("APP_KEY")))
}