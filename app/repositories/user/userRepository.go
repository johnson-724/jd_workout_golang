package user

import (
	"gorm.io/gorm"
	"jd_workout_golang/app/models"
	db "jd_workout_golang/lib/database"
)

func GetUserByEmail(email string) (*models.User, error) {
	user := models.User{}

	result := db.Connection.Where("email = ?", email).Find(&user)

	if result.Error != nil && result.Error == gorm.ErrRecordNotFound {
		return nil, result.Error
	}

	return &user, nil
}

func Update(user *models.User) error {
	result := db.Connection.Save(user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func Create(user *models.User) error {
	result := db.Connection.Create(user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
