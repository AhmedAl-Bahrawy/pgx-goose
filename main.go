package main

import (
	"fmt"
	"log"
	"os"

	"github.com/AhmedAl-Bahrawy/pgx-goose/db"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
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
		log.Fatal(err)
		return nil, err
	}

	defer pool.Close()

	return pool, nil
}

func main() {
	databaseURL, _ := loadEnv()
	pool, _ := loadPool(databaseURL)
	fmt.Printf("Successfully connected to PostgreSQL database. Connection pool: %p\n", pool)

	fmt.Println("First App Lunch\n")
}
