package main

import (
	"jd_workout_golang/lib/database"
)

type user struct {
	EmailVerified int16 `json:"emailVerified" gorm:"default:0"`
}

func UpAddEmailVerifiedColumn() {
	database.Connection.Migrator().AddColumn(&user{}, "EmailVerified")
}

func DownAddEmailVerifiedColumn() {
	database.Connection.Migrator().DropColumn(&user{}, "EmailVerified")
}