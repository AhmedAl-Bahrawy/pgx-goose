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

	return pool, nil
}

func HashPassword(password string) (string, error) {
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
	// var users *[]handlers.User
	//
	var id int64 = 4
	var newUsername string = "hamed"
	var newPasswrod string = "@123AAA"
	//
	// test field

	newPasswordHash, err := HashPassword(newPasswrod)
	if err != nil {
		log.Fatal(err)
	}

	user, err = handlers.UpdateUserById(ctx, pool, id, newUsername, newPasswordHash)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Updated User new username: %v, new password hash: %v", user.Username, user.PasswordHash)
}
