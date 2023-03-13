package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitDatabase 初始化資料庫連線
// 後續改 singlton
func InitDatabase() *gorm.DB {
	db, err := gorm.Open(
		mysql.New(mysql.Config{
			DSN: "root:root@tcp(mysql:3306)/jd_workout?charset=utf8mb4&parseTime=True&loc=Local",
		}),
		&gorm.Config{},
	)

	if err != nil {
		panic(err)
	}

	return db
}