package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"example.com/auth_service/handlers"
	"example.com/auth_service/repository/postgres"

	_ "example.com/auth_service/repository"
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
	log.Printf("Success connection to dbname:'%s', from %s:%s",
		os.Getenv("DB_NAME"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
	)
	defer db.Close()

	postgresRepo := postgres.NewRefreshTokenRepository(db)

	h := handlers.NewHandler(postgresRepo, os.Getenv("JWT_SECRET_KEY"))

	http.HandleFunc("/token", h.HandleGenerateTokens)
	http.HandleFunc("/refresh", h.HandleUpdateTokens)

	log.Printf("Start listen on %s", "8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
