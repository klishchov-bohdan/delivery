package mysql

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

type MySQL struct {
	*gorm.DB
}

// Dial connect to DB and return connection
func Dial() (*MySQL, error) {
	err := godotenv.Load("config/mysql.env")
	if err != nil {
		return nil, err
	}
	mysqlUser := os.Getenv("user")
	mysqlPassword := os.Getenv("pwd")
	mysqlDBName := os.Getenv("dbName")

	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysqlUser, mysqlPassword, mysqlDBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &MySQL{db}, err
}
