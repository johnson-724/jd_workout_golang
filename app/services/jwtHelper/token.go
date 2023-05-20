package jwtHelper

import (
	"fmt"
	"jd_workout_golang/app/models"
	"os"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	env "github.com/joho/godotenv"
)

func GenerateTokenWithPayload(u *models.User, payload map[string]interface{}) (string, error) {
	claims := jwt.MapClaims{
		"uid": u.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	for k, v := range payload {
		claims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)

	return token.SignedString([]byte(os.Getenv("APP_KEY")))
}

func GenerateToken(u *models.User) (string, error) {
	env.Load()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":          u.ID,
		"restPassword": u.ResetPassword,
		"exp":          time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("APP_KEY")))
}

func ValidateToken(tokenString string, uid *uint) (string, bool) {
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

	if int16(claims["restPassword"].(float64)) == 1 {
		return fmt.Errorf("請先重置密碼").Error(), false
	}

	println(uint(claims["uid"].(float64)))

	uidPayload, _ := claims["uid"].(float64)

	*uid = uint(uidPayload)

	return "", true
}

func parseToken(tokenString string) (string, bool) {

	tokenMap := strings.Split(tokenString, "Bearer ")

	if len(tokenMap) != 2 {
		return "", false
	}

	return tokenMap[1], true
}
