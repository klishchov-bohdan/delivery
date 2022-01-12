package mysql

import (
	"errors"
	"github.com/klishchov-bohdan/delivery/internal/models"
)

type UserMySQL struct {
	db *MySQL
}

func NewUserMySQL(db *MySQL) *UserMySQL {
	return &UserMySQL{db: db}
}

func (userMySQL *UserMySQL) GetAll() (users *[]models.UserDB, err error) {
	// select * from `users`
	result := userMySQL.db.Table("users").Find(&users)
	err = result.Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (userMySQL *UserMySQL) Create(user *models.UserDB) (err error) {
	if user == nil {
		return errors.New("no user provided")
	}
	// INSERT INTO `users` (name, email, pwd_hash) VALUES ("Name", "Email", "PasswordHash")
	err = userMySQL.db.
		Table("users").
		Select("Name", "Email", "PasswordHash").
		Create(&user).Error
	return err
}
