package models

import (
	"gorm.io/gorm"
)

type Equip struct {
	gorm.Model
	UserId  uint
	Name    string
	Weights string `gorm:"default:null"`
	Note    string `gorm:"default:null"`
}
