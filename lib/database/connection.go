package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var Connection *gorm.DB

// InitDatabase 初始化資料庫連線
// 後續改 singlton
func InitDatabase() {
	config := gorm.Config{}
	if os.Getenv("LOG_QUERY") == "true" {
		config = gorm.Config{
			Logger: logger.New(
				log.New(os.Stdout, "\r\n", log.LstdFlags),
				logger.Config{
					SlowThreshold: time.Second,
					LogLevel:      logger.Info,
					Colorful:      true,
				},
			),
		}
	}

	db, err := gorm.Open(
		mysql.New(mysql.Config{
			DSN: os.Getenv("DB_HOST"),
		}),
		&config,
	)

	if err != nil {
		panic(err)
	}

	Connection = db
}
