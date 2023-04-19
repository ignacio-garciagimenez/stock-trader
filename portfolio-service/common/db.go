package common

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB() (*gorm.DB, error) {
	//TODO: config env variables
	connectionString := "root:root@tcp(127.0.0.1:3306)/portfolio?charset=utf8mb4&parseTime=True&collation=utf8mb4_0900_ai_ci"
	return gorm.Open(mysql.Open(connectionString), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger: logger.Default.LogMode(logger.Info),
	})
}
