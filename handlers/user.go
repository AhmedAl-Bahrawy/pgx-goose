package handlers

import (
	"context"
	"errors"
	"fmt"
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

func GetUserById(ctx context.Context, pool *pgxpool.Pool, id int64) (*User, error) {
	query := "SELECT id, username, password_hash, created_at FROM users WHERE id = $1"
	row := pool.QueryRow(ctx, query, id)

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

func DeleteUserByUsername(ctx context.Context, pool *pgxpool.Pool, username string) error {
	query := "SELECT id, username FROM users WHERE username = $1"
	row := pool.QueryRow(ctx, query, username)

	var user User
	err := row.Scan(&user.ID, &user.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			fmt.Printf("no rows found in users for username: %v\n", username)
			return nil
		}
		return err
	}

	query = "DELETE FROM users WHERE username = $1"
	tag, err := pool.Exec(ctx, query, username)
	if err != nil {
		return err
	}
	fmt.Printf("Deleted user where is id: %v and username is: %v\n", user.ID, user.Username)
	fmt.Printf("Deleted %d row(s)\n", tag.RowsAffected())
	return nil
}

func DeleteUserById(ctx context.Context, pool *pgxpool.Pool, id int64) error {
	query := "SELECT id, username FROM users WHERE id = $1"
	row := pool.QueryRow(ctx, query, id)

	var user User
	err := row.Scan(&user.ID, &user.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			fmt.Printf("no rows found in users for id: %v\n", id)
			return nil
		}
		return err
	}

	query = "DELETE FROM users WHERE id = $1"
	tag, err := pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	fmt.Printf("Deleted user where is id: %v and username is: %v\n", user.ID, user.Username)
	fmt.Printf("Deleted %d row(s)\n", tag.RowsAffected())
	return nil
}

// Implemented Functions
// 1st. Get All Users
// 2nd. Update User By Id
// 3rd. Update User By Username

// Done Functions
// None Yet
