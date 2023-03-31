package auth

import (
	"gorm.io/gorm"
)

type user struct {
	gorm.Model
	Username string
	Email    string `gorm:"unique"`
	Password string `gorm:"size:64"`
}