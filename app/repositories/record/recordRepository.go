package record

import (
	"fmt"
	"jd_workout_golang/app/models"
	pageinate "jd_workout_golang/app/repositories/pageinate"
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

func Delete(record *models.Record) error {
	result := db.Connection.Delete(record)

	if result.Error != nil {
		return result.Error
	}

	return nil
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

func GetRecords(page pageinate.PaginateCondition, uid uint) (*[]models.Record, *int64, error) {
	data := []models.Record{}
	count := int64(0)

	query := db.Connection.Model(models.Record{}).Where("user_id = ?", uid)

	query.Count(&count)

	result := query.Preload("Equips").Order("created_at desc").Scopes(pageinate.Paginate(page.Page, page.PerPage)).Find(&data)	

	if result.Error != nil {
		return nil, &count, result.Error
	}

	return &data, &count, nil
}