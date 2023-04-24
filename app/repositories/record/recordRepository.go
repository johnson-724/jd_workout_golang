package record

import (
	"fmt"
	"jd_workout_golang/app/models"
	db "jd_workout_golang/lib/database"
	"gorm.io/gorm"
)

func Create(record *models.Record) (uint, error) {
	result := db.Connection.Create(record)

	if result.Error != nil {
		return 0, result.Error
	}

	return record.ID, nil
}

func Update(record *models.Record) error {
	result := db.Connection.Save(record)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func GetRecord(recordId uint64, uid uint) (*models.Record, error) {
	record := models.Record{}

	result := db.Connection.Where("user_id = ?", uid).Where("id = ? ", recordId).First(&record)

	if result.Error != nil && result.Error == gorm.ErrRecordNotFound {
		return &models.Record{}, fmt.Errorf("record not found : %w", result.Error)
	}

	return &record, nil
}