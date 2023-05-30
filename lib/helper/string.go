package helper

import (
	"math/rand"
	"time"
)

func RandString(length int) string {
	letters := []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, length)
    rand.Seed(time.Now().UnixNano())
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }

    return string(b)
}