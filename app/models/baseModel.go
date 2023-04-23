package models

import (
	"time"

	"gorm.io/gorm"
)

// temporarily mark DeletedAt as disabled on JSON conversion, as swaggo cannot parse third-party packages
type baseModel struct {
	ID        uint `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}