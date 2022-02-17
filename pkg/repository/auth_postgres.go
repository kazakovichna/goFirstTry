package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/kazakovichna/todoListPrjct"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) SetSessions(username, refreshToken string, expiresAt int) error {
	query := fmt.Sprintf("UPDATE %s ut SET refreshToken=$1, expiresAt=$2 WHERE ut.username=$3", userTable)
	_, err := r.db.Exec(query, refreshToken, expiresAt, username)
	if err != nil {
		return err
	}

	return nil
}

func (r *AuthPostgres) CreateUser(user todoListPrjct.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, username, passwordHash) values ($1, $2, $3) RETURNING userId", userTable)
	row := r.db.QueryRow(query, user.Name, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *AuthPostgres) GetUser(username, password string) (todoListPrjct.User, error) {
	var user todoListPrjct.User
	query := fmt.Sprintf("SELECT userId FROM %s WHERE username=$1 and passwordHash=$2", userTable)
	err := r.db.Get(&user, query, username, password)

	return user, err
}

func (r *AuthPostgres) GetUserByRefreshToken(refreshToken string) (todoListPrjct.UserRefreshToken, error) {
	var userName todoListPrjct.UserRefreshToken

	query := fmt.Sprintf("SELECT * FROM %s WHERE refreshToken=$1", userTable)
	err := r.db.Get(&userName, query, refreshToken)

	return userName, err
}