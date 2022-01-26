package repo

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/klishchov-bohdan/delivery/internal/models"
)

type TokensRepo struct {
	DB *sql.DB
	TX *sql.Tx
}

func NewTokensRepo(db *sql.DB) *TokensRepo {
	return &TokensRepo{DB: db}
}

func (r *TokensRepo) GetAllTokens() (*[]models.Token, error) {
	var tokens []models.Token
	rows, err := r.DB.Query("SELECT id, user_id, access_hash, refresh_hash, created_at, updated_at FROM tokens")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var token models.Token
		err = rows.Scan(&token.ID, &token.UserID, &token.AccessHash, &token.RefreshHash, &token.CreatedAt, &token.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}
	return &tokens, nil
}

func (r *TokensRepo) GetTokenByID(id uuid.UUID) (*models.Token, error) {
	var token models.Token
	uid, err := id.MarshalBinary()
	if err != nil {
		return nil, err
	}
	err = r.DB.QueryRow(
		"SELECT id, user_id, access_hash, refresh_hash, created_at, updated_at FROM tokens WHERE id = ?", uid).
		Scan(&token.ID, &token.UserID, &token.AccessHash, &token.RefreshHash, &token.CreatedAt, &token.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *TokensRepo) GetTokenByUserID(id uuid.UUID) (*models.Token, error) {
	var token models.Token
	uid, err := id.MarshalBinary()
	if err != nil {
		return nil, err
	}
	err = r.DB.QueryRow(
		"SELECT id, user_id, access_hash, refresh_hash, created_at, updated_at FROM tokens WHERE user_id = ?", uid).
		Scan(&token.ID, &token.UserID, &token.AccessHash, &token.RefreshHash, &token.CreatedAt, &token.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *TokensRepo) CreateToken(token *models.Token) (uuid.UUID, error) {
	if token == nil {
		return uuid.Nil, errors.New("no token provided")
	}
	uid, err := token.ID.MarshalBinary()
	if err != nil {
		return uuid.Nil, err
	}
	userID, err := token.UserID.MarshalBinary()
	if err != nil {
		return uuid.Nil, err
	}
	if r.TX != nil {
		stmt, err := r.TX.Prepare("INSERT INTO tokens(id, user_id, access_hash, refresh_hash) VALUES(?, ?, ?, ?)")
		if err != nil {
			return uuid.Nil, err
		}
		_, err = stmt.Exec(uid, userID, token.AccessHash, token.RefreshHash)
		if err != nil {
			return uuid.Nil, err
		}
		return token.ID, nil
	}
	stmt, err := r.DB.Prepare("INSERT INTO tokens(id, user_id, access_hash, refresh_hash) VALUES(?, ?, ?, ?)")
	if err != nil {
		return uuid.Nil, err
	}
	_, err = stmt.Exec(uid, userID, token.AccessHash, token.RefreshHash)
	if err != nil {
		return uuid.Nil, err
	}
	return token.ID, nil
}

func (r *TokensRepo) UpdateToken(token *models.Token) (uuid.UUID, error) {
	if token == nil {
		return uuid.Nil, errors.New("no token provided")
	}
	uid, err := token.ID.MarshalBinary()
	if err != nil {
		return uuid.Nil, err
	}
	userID, err := token.UserID.MarshalBinary()
	if err != nil {
		return uuid.Nil, err
	}
	if r.TX != nil {
		stmt, err := r.TX.Prepare("UPDATE tokens SET user_id = ?, access_hash = ?, refresh_hash = ? WHERE id = ?")
		if err != nil {
			return uuid.Nil, err
		}
		_, err = stmt.Exec(userID, token.AccessHash, token.RefreshHash, uid)
		if err != nil {
			return uuid.Nil, err
		}
		return token.ID, nil
	}
	stmt, err := r.DB.Prepare("UPDATE tokens SET user_id = ?, access_hash = ?, refresh_hash = ? WHERE id = ?")
	if err != nil {
		return uuid.Nil, err
	}
	_, err = stmt.Exec(userID, token.AccessHash, token.RefreshHash, uid)
	if err != nil {
		return uuid.Nil, err
	}
	return token.ID, nil
}

func (r *TokensRepo) DeleteTokenByID(id uuid.UUID) (uuid.UUID, error) {
	uid, err := id.MarshalBinary()
	if err != nil {
		return uuid.Nil, err
	}
	if r.TX != nil {
		_, err = r.TX.Exec("DELETE FROM tokens WHERE id = ?", uid)
		if err != nil {
			return uuid.Nil, err
		}
		return id, nil
	}
	_, err = r.DB.Exec("DELETE FROM tokens WHERE id = ?", uid)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (r *TokensRepo) DeleteTokenByUserID(id uuid.UUID) (uuid.UUID, error) {
	uid, err := id.MarshalBinary()
	if err != nil {
		return uuid.Nil, err
	}
	var tokenID uuid.UUID
	err = r.DB.QueryRow("SELECT id FROM tokens WHERE user_id = ?", uid).Scan(&tokenID)
	if err != nil {
		return uuid.Nil, err
	}
	if r.TX != nil {
		_, err = r.TX.Exec("DELETE FROM tokens WHERE user_id = ?", uid)
		if err != nil {
			return uuid.Nil, err
		}
		return tokenID, nil
	}
	_, err = r.DB.Exec("DELETE FROM tokens WHERE user_id = ?", uid)
	if err != nil {
		return uuid.Nil, err
	}
	return tokenID, nil
}

func (r *TokensRepo) BeginTx() error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	r.TX = tx
	return nil
}

func (r *TokensRepo) CommitTx() error {
	defer func() {
		r.TX = nil
	}()
	if r.TX != nil {
		return r.TX.Commit()
	}
	return nil
}

func (r *TokensRepo) RollbackTx() error {
	defer func() {
		r.TX = nil
	}()
	if r.TX != nil {
		return r.TX.Rollback()
	}
	return nil
}

func (r *TokensRepo) GetTx() *sql.Tx {
	return r.TX
}

func (r *TokensRepo) SetTx(tx *sql.Tx) {
	r.TX = tx
}
