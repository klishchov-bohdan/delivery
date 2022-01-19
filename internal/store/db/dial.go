package db

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"time"
)

func Dial() (*sql.DB, error) {
	err := godotenv.Load("config/mysql.env")
	if err != nil {
		return nil, err
	}
	mysqlUser := os.Getenv("user")
	mysqlPassword := os.Getenv("pwd")
	mysqlDBName := os.Getenv("dbName")

	db, err := sql.Open(
		"mysql",
		fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?parseTime=true", mysqlUser, mysqlPassword, mysqlDBName),
	)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	maxIdleConns, err := strconv.Atoi(os.Getenv("maxIdleConns"))
	maxOpenConns, err := strconv.Atoi(os.Getenv("maxOpenConns"))
	connMaxLifetime, err := strconv.Atoi(os.Getenv("connMaxLifetime"))
	if err != nil {
		return nil, errors.New("cant parse connecting pool variables")
	}
	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Second)
	return db, nil
}
