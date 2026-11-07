# pgx-goose

A learning project where I explore **goose** (SQL migrations) and **pgx v5** (PostgreSQL driver for Go), and how to use them together — connection pooling, migrations, and basic CRUD queries.

## Project structure

```
.
├── main.go                 # entry point: env, DB connection, demo create/fetch user
├── db/
│   └── postgres.go         # ConnectDB: pgxpool config, connection, ping
├── handlers/
│   └── user.go             # User model + CreateUser / GetUserByUsername
└── sql/
    └── migrations/         # goose migration files
```

## What's inside

- `main.go` — loads env vars from a `.env` file with [godotenv](https://github.com/joho/godotenv), connects to PostgreSQL, then demos an end-to-end flow: hash a password with **bcrypt**, create a user, and fetch it back by username.
- `db/postgres.go` — `ConnectDB()` builds pool config with `pgxpool.ParseConfig`, creates a `pgxpool.Pool`, and pings the database to verify the connection.
- `handlers/user.go` — a `User` model and two query functions:
  - `CreateUser` — `INSERT ... RETURNING` to insert a user and get the full row back.
  - `GetUserByUsername` — `SELECT` by username, returning `nil` cleanly when no rows match (`pgx.ErrNoRows`).
- `sql/migrations/` — goose migrations with `-- +goose Up` / `-- +goose Down` sections:
  1. `init_db` — initial setup
  2. `create_tasks` — `tasks` table (id, title, completed, created_at)
  3. `user_task_user_id` — `users` table + `user_id` foreign key on `tasks` (with `ON DELETE CASCADE`)
  4. `add_password_field_to_user` — adds `password_hash` column
  5. `fix_password_hash_length` — widens `password_hash` to `VARCHAR(255)` (bcrypt hashes are 60 chars)

## Database schema

- **users** — `id`, `username` (unique), `password_hash`, `created_at`
- **tasks** — `id`, `title`, `completed`, `created_at`, `user_id` → users(id)

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

Go `1.27.3`
