package handlers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
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

func CreateTask(ctx context.Context, pool *pgxpool.Pool, user *User, title string, completed bool) (*Task, error) {
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
	query := `
		SELECT
			t.id,
			t.user_id,
			t.title,
			t.completed,
			t.created_at,
			u.id,
			u.username,
			u.password_hash,
			u.created_at
		FROM tasks t
		LEFT JOIN users u ON t.user_id = u.id
		WHERE t.id = $1
	`

	var task Task
	var user User

	err := pool.QueryRow(ctx, query, id).Scan(
		&task.ID,
		&task.UserID,
		&task.Title,
		&task.Completed,
		&task.CreatedAt,
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if task.UserID != nil {
		task.User = &user
	}

	return &task, nil
}

func GetTasksByUser(ctx context.Context, pool *pgxpool.Pool, user *User) (*[]Task, error) {
	query := `
		SELECT
			t.id,
			t.user_id,
			t.title,
			t.completed,
			t.created_at,
			u.id,
			u.username,
			u.created_at
		FROM tasks t
		INNER JOIN users u ON t.user_id = u.id
		WHERE u.id = $1;
	`

	rows, err := pool.Query(ctx, query, user.ID)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	tasks, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (Task, error) {
		var task Task
		var user User

		err := row.Scan(
			&task.ID,
			&task.UserID,
			&task.Title,
			&task.Completed,
			&task.CreatedAt,
			&user.ID,
			&user.Username,
			&user.CreatedAt,
		)
		if err != nil {
			return Task{}, err
		}

		task.User = &user

		return task, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to collect rows: %w", err)
	}

	return &tasks, nil
}

func GetAllTasks(ctx context.Context, pool *pgxpool.Pool) (*[]Task, error) {
	query := `
		SELECT
			t.id AS id,
			t.user_id AS user_id,
			t.title AS title,
			t.completed AS completed,
			t.created_at AS created_at,
			u.id AS user_id_join,
			u.username AS user_username,
			u.created_at AS user_created_at
		FROM tasks t
		LEFT JOIN users u ON t.user_id = u.id
	`

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	tasks, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (Task, error) {
		var task Task
		var user User

		err := row.Scan(
			&task.ID,
			&task.UserID,
			&task.Title,
			&task.Completed,
			&task.CreatedAt,
			&user.ID,
			&user.Username,
			&user.CreatedAt,
		)
		if err != nil {
			return Task{}, err
		}

		if task.UserID != nil {
			task.User = &user
		}

		return task, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to collect tasks: %w", err)
	}

	return &tasks, nil
}

func UpdateTaskById(ctx context.Context, pool *pgxpool.Pool, id int64, newTitle string, newCompleted bool) (*Task, error) {
	query := `
		UPDATE tasks
		SET title = $1, completed = $2
		WHERE id = $3
		RETURNING id, user_id, title, completed, created_at`

	var task Task
	err := pool.QueryRow(ctx, query, newTitle, newCompleted, id).Scan(
		&task.ID,
		&task.UserID,
		&task.Title,
		&task.Completed,
		&task.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("task with id %d not found: %w", id, err)
		}
		return nil, err
	}

	if task.UserID != nil {
		var user User
		err = pool.QueryRow(ctx,
			"SELECT id, username, password_hash, created_at FROM users WHERE id = $1",
			*task.UserID,
		).Scan(
			&user.ID,
			&user.Username,
			&user.PasswordHash,
			&user.CreatedAt,
		)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return &task, nil
			}
			return nil, err
		}
		task.User = &user
	}

	fmt.Printf("Updated task where ID is: %v, new title is: %v and new completed is: %v\n",
		task.ID, task.Title, task.Completed)

	return &task, nil
}

func DeleteTaskById(ctx context.Context, pool *pgxpool.Pool, id int64) error {
	query := "DELETE FROM tasks WHERE id = $1"

	tag, err := pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("task with id %d not found: %w", id, pgx.ErrNoRows)
	}

	fmt.Printf("task with id %d deleted successfully\n", id)

	return nil
}
