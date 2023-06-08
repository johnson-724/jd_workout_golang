package equip

import (
	"fmt"
	"gorm.io/gorm"
	"jd_workout_golang/app/models"
	db "jd_workout_golang/lib/database"
)

type PaginateCondition struct {
	Page    int
	PerPage int
}

func Create(equip *models.Equip) (uint, error) {
	result := db.Connection.Create(equip)

	if result.Error != nil {
		return 0, result.Error
	}

	return equip.ID, nil
}

func Delete(equip *models.Equip) error {
	result := db.Connection.Delete(equip)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func Update(equip *models.Equip) error {
	result := db.Connection.Model(equip).Updates(map[string]interface{}{
		"name":    equip.Name,
		"note":    equip.Note,
		"image":   equip.Image,
		"weights": equip.Weights,
	})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func GetEquip(equipId uint64, uid uint) (*models.Equip, error) {
	equip := models.Equip{}

	result := db.Connection.Where("user_id = ?", uid).Where("id = ? ", equipId).First(&equip)

	if result.Error != nil && result.Error == gorm.ErrRecordNotFound {
		return &models.Equip{}, fmt.Errorf("equip not found : %w", result.Error)
	}

	return &equip, nil
}

func GetEqupis(page PaginateCondition, uid uint) (*[]models.Equip, *int64, error) {
	data := []models.Equip{}
	count := int64(0)

	query := db.Connection.Model(models.Equip{}).Where("user_id = ?", uid)

	query.Count(&count)

	result := query.Order("name asc").Scopes(Paginate(page.Page, page.PerPage)).Find(&data)

	if result.Error != nil {
		return nil, &count, result.Error
	}

	return &data, &count, nil
}

func Paginate(currentPage int, perPage int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		fmt.Printf("currentPage: %d, perPage: %d", currentPage, perPage)
		offset := (currentPage - 1) * perPage

		return db.Offset(offset).Limit(perPage)
	}
}
