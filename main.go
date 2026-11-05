package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func ConnectDB() (*pgxpool.Pool, error) {
	databaseURL := os.Getenv("GOOSE_DBSTRING")

	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}


func main(){
	err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

	pool, err := ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	fmt.Println("First App Lunch")
}