package main

import (
	"jd_workout_golang/app/models"
	"jd_workout_golang/lib/database"
)

func UpAlterRestPasswordColumn() {
	database.Connection.Migrator().AlterColumn(&models.User{}, "ResetPassword")
}

func DownAlterRestPasswordColumn() {
	database.Connection.Migrator().AlterColumn(&models.User{}, "ResetPassword")
}