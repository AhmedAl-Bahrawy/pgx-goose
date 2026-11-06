package main

import (
	"fmt"
	"log"

	"github.com/AhmedAl-Bahrawy/pgx-goose/database"
	"github.com/joho/godotenv"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
func loadPool() {
	pool, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	defer pool.Close()
}

func main() {
	loadEnv()
	loadPool()

	fmt.Println("First App Lunch")
}
