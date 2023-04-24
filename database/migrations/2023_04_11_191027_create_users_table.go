package main

import (
	"jd_workout_golang/lib/database"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"type:varchar(255);unique"`
	Password string `gorm:"type:varchar(255)"`
	Username string `gorm:"type:varchar(255)"`
}

func UpCreateUsersTable() {
	database.Connection.Migrator().CreateTable(&User{})
}

func DownCreateUsersTable() {
	database.Connection.Migrator().DropTable(&User{})
}
