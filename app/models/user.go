package models

type User struct {
	Model
	Username string `json:"username"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password" gorm:"size:64"`
}