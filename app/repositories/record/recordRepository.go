package record

import (
	"jd_workout_golang/app/models"
	db "jd_workout_golang/lib/database"
)

func Create(record *models.Record) (uint, error) {
	result := db.Connection.Create(record)

	if result.Error != nil {
		return 0, result.Error
	}

	return record.ID, nil
}