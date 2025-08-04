package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"example.com/auth_service/config"
	"example.com/auth_service/handlers"
	"example.com/auth_service/repository/postgres"

	_ "example.com/auth_service/repository"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.LoadConfig()

	connStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBHost,
		cfg.DBPort,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Success connection to dbname:'%s', from %s:%s",
		cfg.DBName,
		cfg.DBHost,
		cfg.DBPort,
	)
	defer db.Close()

	postgresRepo := postgres.NewRefreshTokenRepository(db)

	h := handlers.NewHandler(postgresRepo, cfg.JWTSecretKey)

	http.HandleFunc("/token", h.HandleGenerateTokens)
	http.HandleFunc("/refresh", h.HandleUpdateTokens)

	log.Printf("Start listen on %s", "8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
