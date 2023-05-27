package main

import (
	"jd_workout_golang/app/models"
	"jd_workout_golang/lib/database"
)


func UpAddImageToEuqipTable() {
	database.Connection.Migrator().AddColumn(&models.Equip{}, "Image")
}

func DownAddImageToEuqipTable() {
	database.Connection.Migrator().DropColumn(&models.Equip{}, "Image")
}