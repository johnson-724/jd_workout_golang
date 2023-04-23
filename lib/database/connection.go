package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

// InitDatabase 初始化資料庫連線
// 後續改 singlton
func InitDatabase() *gorm.DB {
	db, err := gorm.Open(
		mysql.New(mysql.Config{
			DSN: os.Getenv("DB_HOST"),
		}),
		&gorm.Config{},
	)

	if err != nil {
		panic(err)
	}

	return db
}
