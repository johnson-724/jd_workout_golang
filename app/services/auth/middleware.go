package auth

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	env "github.com/joho/godotenv"
)

var uid float64

func ValidateToken(c *gin.Context) {

	env.Load()
	tokenString := strings.Split(c.GetHeader("Authorization"), "Bearer ")[1]
	// jwt.SigningMethodHS256

	token, _ := jwt.Parse(
		tokenString,
		// func to get the key for validating
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(os.Getenv("APP_KEY")), nil
		})

	// Type Assertion
	// token.Claims is implement jwt.MapClaims type, claims is jwt.MapClaims and ok == true
	// or ok == false and claims is nil
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "JWT token payload is improper",
		})

		c.Abort()

		return
	}

	uid = claims["uid"].(float64)

	fmt.Println(claims["uid"])

	c.Next()
}
