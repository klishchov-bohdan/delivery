package repo

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/klishchov-bohdan/delivery/internal/models"
)

type UsersRepo struct {
	DB *sql.DB
	TX *sql.Tx
}

func NewUsersRepo(db *sql.DB) *UsersRepo {
	return &UsersRepo{DB: db}
}

func (r *UsersRepo) GetAllUsers() (*[]models.User, error) {
	var users []models.User
	rows, err := r.DB.Query("SELECT id, name, email, password_hash, created_at, updated_at FROM users")
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

func (r *UsersRepo) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.DB.QueryRow(
		"SELECT id, name, email, password_hash, created_at, updated_at FROM users WHERE email = ?", email).
		Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UsersRepo) GetUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	uid, err := id.MarshalBinary()
	if err != nil {
		return nil, err
	}
	err = r.DB.QueryRow(
		"SELECT id, name, email, password_hash, created_at, updated_at FROM users WHERE id = ?", uid).
		Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UsersRepo) CreateUser(user *models.User) (id uuid.UUID, err error) {
	if user == nil {
		return id, errors.New("no user provided")
	}
	uid, err := user.ID.MarshalBinary()
	if err != nil {
		return id, err
	}
	if r.TX != nil {
		stmt, err := r.TX.Prepare("INSERT INTO users(id, name, email, password_hash) VALUES(?, ?, ?, ?)  RETURNING id")
		if err != nil {
			return id, err
		}
		err = stmt.QueryRow(uid, user.Name, user.Email, user.PasswordHash).Scan(&id)
		if err != nil {
			return id, err
		}
		return id, nil
	}
	stmt, err := r.DB.Prepare("INSERT INTO users(id, name, email, password_hash) VALUES(?, ?, ?, ?)")
	if err != nil {
		return id, err
	}
	_, err = stmt.Exec(uid, user.Name, user.Email, user.PasswordHash)
	if err != nil {
		return id, err
	}
	return id, nil
}

func (r *UsersRepo) UpdateUser(user *models.User) (err error) {
	if user == nil {
		return errors.New("no user provided")
	}
	err = r.DB.QueryRow("SELECT * FROM users WHERE id = ?", user.ID).Scan()
	if err != nil {
		return errors.New("user not found")
	}
	uid, err := user.ID.MarshalBinary()
	if err != nil {
		return err
	}
	if r.TX != nil {
		stmt, err := r.TX.Prepare("UPDATE users SET name = ?, email = ?, password_hash = ? WHERE id = ?")
		if err != nil {
			return err
		}
		_, err = stmt.Exec(user.Name, user.Email, user.PasswordHash, uid)
		if err != nil {
			return err
		}
		return nil
	}
	stmt, err := r.DB.Prepare("UPDATE users SET name = ?, email = ?, password_hash = ? WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(user.Name, user.Email, user.PasswordHash, uid)
	if err != nil {
		return err
	}
	return nil
}

func (r *UsersRepo) DeleteUser(id uuid.UUID) (err error) {
	err = r.DB.QueryRow("SELECT * FROM users WHERE id = ?", id).Scan()
	if err != nil {
		return errors.New("user not found")
	}
	uid, err := id.MarshalBinary()
	if err != nil {
		return err
	}
	if r.TX != nil {
		_, err = r.TX.Exec("DELETE FROM users WHERE id = ?", uid)
		if err != nil {
			return err
		}
		return nil
	}
	_, err = r.DB.Exec("DELETE FROM users WHERE id = ?", uid)
	if err != nil {
		return err
	}
	return nil
}

func (r *UsersRepo) BeginTx() error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	r.TX = tx
	return nil
}

func (r *UsersRepo) CommitTx() error {
	defer func() {
		r.TX = nil
	}()
	if r.TX != nil {
		return r.TX.Commit()
	}
	return nil
}

func (r *UsersRepo) RollbackTx() error {
	defer func() {
		r.TX = nil
	}()
	if r.TX != nil {
		return r.TX.Rollback()
	}
	return nil
}

func (r *UsersRepo) GetTx() *sql.Tx {
	return r.TX
}

func (r *UsersRepo) SetTx(tx *sql.Tx) {
	r.TX = tx
}
