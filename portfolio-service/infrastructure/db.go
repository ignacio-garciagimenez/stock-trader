package infrastructure

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB() (*gorm.DB, error) {
	//TODO: config env variables
	connectionString := fmt.Sprintf("root:root@tcp(%s:3306)/portfolio?charset=utf8mb4&parseTime=True&collation=utf8mb4_0900_ai_ci", os.Getenv("MYSQL_HOST"))
	return gorm.Open(mysql.Open(connectionString), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(logger.Info),
	})
}
