package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/klishchov-bohdan/delivery/config"
	"time"
)

func Dial() (*sql.DB, error) {
	cfg := config.NewConfig()

	db, err := sql.Open(
		"mysql",
		fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?parseTime=true", cfg.MysqlUserName, cfg.MysqlPassword, cfg.MysqlDBName),
	)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(cfg.MysqlMaxIdleCons)
	db.SetMaxOpenConns(cfg.MysqlMaxOpenCons)
	db.SetConnMaxLifetime(time.Duration(cfg.MysqlConsMaxLifeTime) * time.Second)
	return db, nil
}
