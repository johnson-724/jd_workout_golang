package main

import (
	"jd_workout_golang/app/models"
	"jd_workout_golang/lib/database"
)

func UpAddResetPasswordColumnToUserTable() {
	database.Connection.Migrator().AddColumn(&models.User{}, "ResetPassword")
}

func DownAddResetPasswordColumnToUserTable() {
	database.Connection.Migrator().DropColumn(&models.User{}, "ResetPassword")
}