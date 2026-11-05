# pgx-goose

A learning project where I explore **goose** (SQL migrations) and **pgx** (PostgreSQL driver for Go), and how to use them together.

## What's inside

- `main.go` — entry point that:
  - loads environment variables from a `.env` file using [godotenv](https://github.com/joho/godotenv),
  - connects to PostgreSQL with [pgx v5](https://github.com/jackc/pgx) using a connection pool (`pgxpool`),
  - pings the database to verify the connection before use.
- `sql/migrations/` — goose migration files (e.g. `init_db.sql` with the `-- +goose Up` / `-- +goose Down` sections).

## Environment

Requires a `.env` file with the database connection string used by goose/pgx:

```
GOOSE_DBSTRING=postgres://user:password@localhost:5432/dbname
```

## Running

```bash
go run main.go
```

## Go version

Go `1.27.3`
