package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

var Connection *gorm.DB

// InitDatabase 初始化資料庫連線
// 後續改 singlton
func InitDatabase() {
	db, err := gorm.Open(
		mysql.New(mysql.Config{
			DSN: os.Getenv("DB_HOST"),
		}),
		&gorm.Config{},
	)

	if err != nil {
		panic(err)
	}

	Connection = db
}
