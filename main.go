package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/AhmedAl-Bahrawy/pgx-goose/db"
	"github.com/AhmedAl-Bahrawy/pgx-goose/handlers"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func loadEnv() (*string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return nil, err
	}

	databaseURL := os.Getenv("GOOSE_DBSTRING")
	if databaseURL == "" {
		log.Fatal("DB_URL is empty")
	}
	return &databaseURL, nil
}

func loadPool(databaseURL *string) (*pgxpool.Pool, error) {
	pool, err := db.ConnectDB(databaseURL)
	if err != nil {
		return nil, err
	}

	// Note: do NOT defer pool.Close() here, it would close the pool
	// as soon as this function returns, before main ever uses it.

	return pool, nil
}

func HashPassword(password string) (string, error) {
	// GenerateFromPassword automatically manages a secure, random salt internally
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func main() {
	ctx := context.Background()

	databaseURL, err := loadEnv()
	if err != nil {
		log.Fatal(err)
	}

	pool, err := loadPool(databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	fmt.Printf("Successfully connected to PostgreSQL database. Connection pool: %p\n", pool)

	var user *handlers.User
	username := "ahmedalbahrawy"
	password := "123456"

	password_hash, err := HashPassword(password)
	if err != nil {
		log.Fatal(err)
	}

	user, err = handlers.CreateUser(ctx, pool, username, password_hash)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created user: %+v\n", *user)
}
