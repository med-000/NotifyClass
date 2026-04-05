package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDB() (*gorm.DB, error) {
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	database := os.Getenv("MYSQL_DATABASE")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		user, password, host, port, database,
	)

	newLogger := logger.New(
		log.New(os.Stdout, "", log.LstdFlags),
		logger.Config{
			LogLevel:                  logger.Warn, 
			IgnoreRecordNotFoundError: true,
		},
	)

	return gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
}