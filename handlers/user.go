package handlers

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID           int64     `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	PasswordHash string    `json:"password_hash" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

func CreateUser(ctx context.Context, pool *pgxpool.Pool, username string, password_hash string) (*User, error) {
	// Create User then check no issues
	query := "INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id, username, password_hash, created_at"
	row := pool.QueryRow(ctx, query, username, password_hash)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByUsername(ctx context.Context, pool *pgxpool.Pool, username string) (*User, error) {
	query := "SELECT id, username, password_hash, created_at FROM users WHERE username = $1"
	row := pool.QueryRow(ctx, query, username)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
