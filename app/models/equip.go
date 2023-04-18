package models

import (
	"gorm.io/gorm"
)

type Equip struct {
	gorm.Model
	UserId  uint
	Name    string
	Weights string
	Note    string
}
