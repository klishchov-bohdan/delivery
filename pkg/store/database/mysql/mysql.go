package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"os"
)

type MySQL struct {
	*sql.DB
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

	db, err := sql.Open(
		"mysql",
		fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?parseTime=true", mysqlUser, mysqlPassword, mysqlDBName),
	)
	if err != nil {
		return nil, err
	}
	return &MySQL{db}, nil
}

func (mysqlConn *MySQL) GetDB() *sql.DB {
	return mysqlConn.DB
}
