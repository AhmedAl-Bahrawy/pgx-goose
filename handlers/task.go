package handlers

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Task struct {
	ID        int64     `json:"id" db:"id"`
	UserID    *int64    `json:"user_id" db:"user_id"`
	Title     string    `json:"title" db:"title"`
	Completed bool      `json:"completed" db:"completed"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`

	User *User `json:"user,omitempty" db:"-"`
}

func CreateTask(
	ctx context.Context,
	pool *pgxpool.Pool,
	user *User,
	title string,
	completed bool,
) (*Task, error) {
	query := `
        INSERT INTO tasks (user_id, title, completed)
        VALUES ($1, $2, $3)
        RETURNING id, user_id, title, completed, created_at
    `

	var task Task

	err := pool.QueryRow(
		ctx, query, user.ID, title, completed,
	).Scan(
		&task.ID,
		&task.UserID,
		&task.Title,
		&task.Completed,
		&task.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	task.User = user

	return &task, nil
}

func GetTaskById(ctx context.Context, pool *pgxpool.Pool, id int64) (*Task, error) {
	query := `SELECT
				id, user_id, title, completed, created_at
				FROM tasks
				WHERE id = $1`

	var task Task
	var user User

	err := pool.QueryRow(ctx, query, id).Scan(
		&task.ID,
		&task.UserID,
		&task.Title,
		&task.Completed,
		&task.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	query = `SELECT
		id, username, password_hash, created_at
		FROM users
		WHERE id = $1`

	err = pool.QueryRow(ctx, query, task.UserID).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	task.User = &user

	return &task, nil
}
