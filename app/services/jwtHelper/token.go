package jwtHelper

import (
	jwt "github.com/golang-jwt/jwt/v5"
	env "github.com/joho/godotenv"
	"os"
	"jd_workout_golang/app/models"
	"time"
	"strings"
)

func GenerateToken (u *models.User) (string, error) {
	env.Load()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": u.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("APP_KEY")))
}

func ValidateToken (tokenString string, uid *float64) (string, bool) {
	token, ok := parseToken(tokenString)

	if !ok {
		return "invalid token", false
	}

	jwtToken, err := jwt.Parse(
		token,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, nil
			}

			return []byte(os.Getenv("APP_KEY")), nil
		})
	
	if err != nil {
		return err.Error(), false
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)

	if !ok || !jwtToken.Valid {
		return err.Error(), false
	}
	
	uidPayload, _ := claims["uid"].(float64)

	*uid = uidPayload
	
	return "", true
}

func parseToken (tokenString string) (string, bool) {

	tokenMap := strings.Split(tokenString, "Bearer ")

	if (len(tokenMap) != 2) {
		return "", false
	}

	return tokenMap[1], true
}