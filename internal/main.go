package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"example.com/auth_service/handlers"
	"example.com/auth_service/db/postgres"

	_ "example.com/auth_service/db"
	_ "github.com/lib/pq"
)

func main() {
	connStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	posresRepo := postgres.NewRefreshTokenRepository(db)

	posresRepo.CheckRefreshToken("123", "e2fsqf")

	http.HandleFunc("/token", handlers.HandleGenerateTokens)
	http.HandleFunc("/refresh", handlers.HandleUpdateTokens)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
