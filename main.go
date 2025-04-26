package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"example.com/auth_service/internal/handlers"

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

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT id, user_id, client_ip FROM refresh_tokens LIMIT 1")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	if rows.Next() {
		var id string
		var userID string
		var clientIP string

		err = rows.Scan(&id, &userID, &clientIP)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("First Row:")
		fmt.Println("ID:", id)
		fmt.Println("UserID:", userID)
		fmt.Println("ClientIP:", clientIP)
	} else {
		fmt.Println("No rows found!")
	}

	columns, err := rows.Columns()
	fmt.Println(columns)

	http.HandleFunc("/token", handler.HandleGenerateTokens)
	http.HandleFunc("/refresh", handler.HandleUpdateTokens)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
