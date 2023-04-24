package models

import (
	"time"

	"gorm.io/gorm"
)

type Record struct {
	ID        uint `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	UserId  uint  `json:"userId"`
	EquipId uint `json:"equip_id"`
	Equip Equip
	Name string `json:"name"`
	Weight float32 `json:"weight" gorm:"default:null"`
	Reps uint `json:"reps" gorm:"default:null"`
	Note string `json:"note" gorm:"default:null"`
}
