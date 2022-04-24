package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Config struct {
	MysqlUserName              string
	MysqlPassword              string
	MysqlDBName                string
	MysqlMaxIdleCons           int
	MysqlMaxOpenCons           int
	MysqlConsMaxLifeTime       int
	AccessTokenLifeTime        int
	RefreshTokenLifeTime       int
	AccessSecret               string
	RefreshSecret              string
	AccessControlAllowOrigin   string
	AccessControlAllowHeaders  string
	AccessControlExposeHeaders string
	AccessControlAllowMethods  string
}

func NewConfig() *Config {
	cfg := &Config{}
	err := godotenv.Load("config/mysql.env")
	if err != nil {
		log.Fatal("Can`t load mysql.env")
	}
	cfg.MysqlUserName = os.Getenv("user")
	cfg.MysqlPassword = os.Getenv("pwd")
	cfg.MysqlDBName = os.Getenv("dbName")
	cfg.MysqlMaxIdleCons, _ = strconv.Atoi(os.Getenv("maxIdleConns"))
	cfg.MysqlMaxOpenCons, _ = strconv.Atoi(os.Getenv("maxOpenConns"))
	cfg.MysqlConsMaxLifeTime, _ = strconv.Atoi(os.Getenv("connMaxLifetime"))

	err = godotenv.Load("config/token.env")
	if err != nil {
		log.Fatal("Cant load token.env")
	}
	cfg.AccessTokenLifeTime, _ = strconv.Atoi(os.Getenv("AccessTokenLifeTime"))
	cfg.RefreshTokenLifeTime, _ = strconv.Atoi(os.Getenv("RefreshTokenLifeTime"))
	cfg.AccessSecret = os.Getenv("AccessSecret")
	cfg.RefreshSecret = os.Getenv("RefreshSecret")

	err = godotenv.Load("config/cors.env")
	if err != nil {
		log.Fatal("Cant load cors.env")
	}

	cfg.AccessControlAllowOrigin = os.Getenv("AccessControlAllowOrigin")
	cfg.AccessControlAllowHeaders = os.Getenv("AccessControlAllowHeaders")
	cfg.AccessControlExposeHeaders = os.Getenv("AccessControlExposeHeaders")
	cfg.AccessControlAllowMethods = os.Getenv("AccessControlAllowMethods")
	return cfg
}
