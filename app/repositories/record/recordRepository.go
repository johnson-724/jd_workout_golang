package record

import (
	"fmt"
	"jd_workout_golang/app/models"
	pageinate "jd_workout_golang/app/repositories/pageinate"
	db "jd_workout_golang/lib/database"
	"time"

	"gorm.io/gorm"
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

func GetRecordList(page pageinate.PaginateCondition, uid uint) (*[]models.Record, *int64, error) {
	records := []models.Record{}
	count := int64(0)

	query := db.Connection.Model(models.Record{}).
		Where("user_id = ?", uid).
		Preload("Equip").
		Scopes(pageinate.Paginate(page.Page, page.PerPage)).
		Find(&records)

	if query.Error != nil {
		return nil, nil, query.Error
	}

	return &records, &count, nil
}

func GetDateSummaryRecords(page pageinate.PaginateCondition, uid uint) (*[]RecordByDate, *int64, error) {
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
		Preload("Equip").
		Where("user_id = ?", uid).
		Where("date_format(created_at, '%Y-%m-%d') IN (?)", format).
		Select("*, date_format(created_at, '%Y-%m-%d') as date").
		Order("created_at desc, equip_id, weight, reps").
		Find(&data)

	return &data, &count, nil
}

type RecordWithVolumn struct {
	models.Record
	Ids    string  `json:"ids"`
	Volumn float32 `json:"volumn"`
	Date   string  `json:"date"`
	Count  int     `json:"sets"`
	Notes  string  `json:"notes"`
}

func GetMaxRecord(equips []uint, before string) *[]RecordWithVolumn {
	records := []RecordWithVolumn{}
	maxWeight := db.Connection.Model(models.Record{}).
		Select("equip_id, max(weight) as weight, max(reps) as reps, "+
			"count(1) as count, date_format(created_at, '%Y-%m-%d') as date, "+
			"row_number() over (partition by equip_id order by weight desc, reps desc, count(1)) as row_num").
		Where("equip_id", equips).
		Where("created_at < ?", before).
		Group("equip_id, weight, reps, date_format(created_at, '%Y-%m-%d')").
		Order("weight desc, reps desc")

	db.Connection.Model(models.Record{}).
		Select("records.id, records.equip_id, records.weight, records.reps, tmp.count, (records.weight * records.reps * tmp.count) as volumn").
		Joins("join (?) as tmp on records.equip_id = tmp.equip_id and records.weight = tmp.weight and records.reps = tmp.reps and row_num = 1 ", maxWeight).
		// Where("records.equip_id = tmp.equip_id and records.weight = tmp.weight and records.reps = tmp.reps and row_num = 1").
		Find(&records)

	return &records
}

func GetRecentRecord(equips []uint) *[]RecordWithVolumn {
	records := []RecordWithVolumn{}

	db.Connection.Model(models.Record{}).
		Select("group_concat(records.id) as ids, records.weight, records.reps, records.equip_id, group_concat(records.note) as notes, "+
			// "row_number() over (partition by equip_id order by date_format(created_at, '%Y-%m-%d')) as row_num,"+
			"date_format(created_at, '%Y-%m-%d') as date, count(1) as count",
		).
		Where("records.created_at >= ?", time.Now().AddDate(0, 0, -7).Format("2006-01-02")).
		Group("equip_id, weight, reps, date_format(created_at, '%Y-%m-%d')").
		Order("equip_id asc, date_format(created_at, '%Y-%m-%d') asc").
		Find(&records)

	return &records
}
