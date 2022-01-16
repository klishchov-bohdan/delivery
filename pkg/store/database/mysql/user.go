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

func (userMySQL *UserMySQL) GetAllUsers() (*[]models.User, error) {
	var users []models.User
	rows, err := userMySQL.db.Query("SELECT id, name, email, password_hash, created_at, updated_at FROM users")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return &users, nil
}

func (userMySQL *UserMySQL) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := userMySQL.db.QueryRow(
		"SELECT id, name, email, password_hash, created_at, updated_at FROM users WHERE email = ?", email).
		Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (userMySQL *UserMySQL) GetUserByID(id []byte) (*models.User, error) {
	var user models.User
	err := userMySQL.db.QueryRow(
		"SELECT id, name, email, password_hash, created_at, updated_at FROM users WHERE id = ?", id).
		Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (userMySQL *UserMySQL) CreateUser(user *models.User) (err error) {
	if user == nil {
		return errors.New("no user provided")
	}
	stmt, err := userMySQL.db.Prepare("INSERT INTO users(id, name, email, password_hash) VALUES(?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(user.ID, user.Name, user.Email, user.PasswordHash)
	if err != nil {
		return err
	}
	return nil
}

func (userMySQL *UserMySQL) UpdateUser(user *models.User) (err error) {
	if user == nil {
		return errors.New("no user provided")
	}
	err = userMySQL.db.QueryRow("SELECT * FROM users WHERE id = ?", user.ID).Scan()
	if err != nil {
		return errors.New("user not found")
	}
	stmt, err := userMySQL.db.Prepare("UPDATE users SET name = ?, email = ?, password_hash = ? WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(user.Name, user.Email, user.PasswordHash, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (userMySQL *UserMySQL) DeleteUser(id []byte) (err error) {
	err = userMySQL.db.QueryRow("SELECT * FROM users WHERE id = ?", id).Scan()
	if err != nil {
		return errors.New("user not found")
	}
	_, err = userMySQL.db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}
