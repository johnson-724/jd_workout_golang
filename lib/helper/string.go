package helper

import "math/rand"

func RandString(length int) string {
	letters := []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, length)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }

    return string(b)
}