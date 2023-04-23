package models

type User struct {
	baseModel
	Username string `json:"username"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password" gorm:"size:64"`
}