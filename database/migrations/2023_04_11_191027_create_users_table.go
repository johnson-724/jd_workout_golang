package main

import (
	"jd_workout_golang/lib/database"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"type:varchar(255);unique_index"`
	Password string `gorm:"type:varchar(255)"`
	Username string `gorm:"type:varchar(255);unique_index"`
}

func UpCreateUsersTable() {
	db := database.InitDatabase()

	db.Migrator().CreateTable(&User{})
}

func DownCreateUsersTable() {
	db := database.InitDatabase()

	db.Migrator().DropTable(&User{})
}
