package record

import (
	"fmt"
	"gorm.io/gorm"
	"jd_workout_golang/app/models"
	pageinate "jd_workout_golang/app/repositories/pageinate"
	db "jd_workout_golang/lib/database"
)

type RecordByDate struct {
	Date string `json:"date"`
	models.Record
}

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

func GetRecords(page pageinate.PaginateCondition, uid uint) (*[]RecordByDate, *int64, error) {
	data := []RecordByDate{}
	groupBy := []struct {
		Date string `json:"date"`
	}{}

	count := int64(0)

	// base query
	query := db.Connection.Model(models.Record{}).Where("user_id = ?", uid)
	// group by date
	groupByQuery := query.Select("date_format(created_at, '%Y-%m-%d') as date ").Group("date").Order("date desc")
	db.Connection.Table("(?) as tmp", groupByQuery).Count(&count)

	groupByQuery = groupByQuery.Scopes(pageinate.Paginate(page.Page, page.PerPage)).Find(&groupBy)
	// run limit in subquery
	format := db.Connection.Table("(?) as tmp", groupByQuery)

	db.Connection.Model(models.Record{}).
		Where("user_id = ?", uid).
		Where("date_format(created_at, '%Y-%m-%d') IN (?)", format).
		Select("*, date_format(created_at, '%Y-%m-%d') as date").
		Order("created_at desc, equip_id, weight, reps").
		Find(&data)

	return &data, &count, nil
}
