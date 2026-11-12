# pgx-goose

A learning project where I explore **goose** (SQL migrations) and **pgx v5** (PostgreSQL driver for Go), and how to use them together — connection pooling, migrations, and basic CRUD queries.

## Project structure

```
.
├── main.go                 # entry point: env loading + DB connection
├── db/
│   └── postgres.go         # ConnectDB: pgxpool config, connection, ping
├── handlers/
│   ├── user.go             # User model + full CRUD
│   └── task.go             # Task model + full CRUD (with user joins)
└── sql/
    └── migrations/         # goose migration files
```

## What's inside

- `main.go` — loads `GOOSE_DBSTRING` from a `.env` file with [godotenv](https://github.com/joho/godotenv), connects to PostgreSQL, and (in a commented-out demo section) exercises the handlers end to end: hashing passwords with **bcrypt**, creating/fetching users and tasks. Also exports a `HashPassword` helper using bcrypt's `DefaultCost`.
- `db/postgres.go` — `ConnectDB()` builds pool config with `pgxpool.ParseConfig`, creates a `pgxpool.Pool` via `NewWithConfig`, and pings the database to verify the connection (closing the pool on failure).
- `handlers/user.go` — a `User` model with JSON/`db` tags and full CRUD:
  - `CreateUser` — `INSERT ... RETURNING` to insert a user and get the full row back.
  - `GetUserByUsername` / `GetUserById` — `SELECT` by key, returning `nil` cleanly when no rows match (`pgx.ErrNoRows`).
  - `GetAllUsers` — lists all users using `pgx.CollectRows` with `pgx.RowToStructByName[User]`.
  - `UpdateUserById` / `UpdateUserByUsername` — check existence first, then `UPDATE ... RETURNING` the new row.
  - `DeleteUserByUsername` / `DeleteUserById` — look up the user first, then `DELETE` and report `RowsAffected`.
- `handlers/task.go` — a `Task` model (with an optional embedded `*User`) and full CRUD:
  - `CreateTask` — inserts a task for a given user, returns the full row.
  - `GetTaskById` / `GetAllTasks` — `LEFT JOIN users` to eager-load the task's owner; `GetTasksByUser` uses an `INNER JOIN` and collects rows with a custom scan function.
  - `UpdateTaskById` — `UPDATE ... RETURNING`, then fetches the owning user separately.
  - `DeleteTaskById` — `Exec` + `RowsAffected`, returning a "not found" error when nothing was deleted.
- `sql/migrations/` — goose migrations with `-- +goose Up` / `-- +goose Down` sections:
  1. `init_db` — initial setup
  2. `create_tasks` — `tasks` table (id, title, completed, created_at)
  3. `user_task_user_id` — `users` table + `user_id` foreign key on `tasks` (with `ON DELETE CASCADE`)
  4. `add_password_field_to_user` — adds `password_hash` column
  5. `fix_password_hash_length` — widens `password_hash` to `VARCHAR(255)` (bcrypt hashes are 60 chars)

## Database schema

- **users** — `id`, `username` (unique), `password_hash`, `created_at`
- **tasks** — `id`, `title`, `completed`, `created_at`, `user_id` → users(id) (nullable in the model, `ON DELETE CASCADE`)

## Models

Both models carry `json` and `db` struct tags:

```go
type User struct {
    ID           int64     `json:"id" db:"id"`
    Username     string    `json:"username" db:"username"`
    PasswordHash string    `json:"-" db:"password_hash"` // hidden from JSON
    CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type Task struct {
    ID        int64     `json:"id" db:"id"`
    UserID    *int64    `json:"user_id" db:"user_id"`
    Title     string    `json:"title" db:"title"`
    Completed bool      `json:"completed" db:"completed"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    User      *User     `json:"user,omitempty" db:"-"` // eager-loaded owner
}
```

## Migrations

Run with [goose](https://pressly.github.io/goose/):

```bash
goose up     # apply all pending migrations
goose down   # roll back the last one
goose status # see migration state
```

## Environment

Requires a `.env` file with the database connection string:

```
GOOSE_DBSTRING=postgres://user:password@localhost:5432/dbname
```

## Running

```bash
go run main.go
```

## Key libraries

- [pgx v5](https://github.com/jackc/pgx) — PostgreSQL driver + connection pool
- [goose v3](https://github.com/pressly/goose) — SQL migrations
- [godotenv](https://github.com/joho/godotenv) — `.env` loading
- [golang.org/x/crypto/bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) — password hashing

## Go version

Go `1.27.3` (per `go.mod`)

