package models

import (
	"time"
	"gorm.io/gorm"
)

type Equip struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	UserId  uint  `json:"userId"`
	Name    string `json:"name"`
	Weights string `json:"weights" gorm:"default:null"`
	Note    string `json:"note" gorm:"default:null"`
}
