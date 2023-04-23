package models

type Equip struct {
	baseModel
	UserId  uint  `json:"userId"`
	Name    string `json:"name"`
	Weights string `json:"weights" gorm:"default:null"`
	Note    string `json:"note" gorm:"default:null"`
}
