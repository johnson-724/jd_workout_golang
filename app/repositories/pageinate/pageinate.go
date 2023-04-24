package pageinate

import (
	"fmt"
	"gorm.io/gorm"
)

type PaginateCondition struct {
	Page    int
	PerPage int
}

func Paginate(currentPage int, perPage int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (currentPage - 1) * perPage

		return db.Offset(offset).Limit(perPage)
	}
}