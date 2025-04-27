package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
	_ "example.com/auth_service/db"
)

type PostgresTokenRepository struct {
	DB *sql.DB
}

func NewRefreshTokenRepository(db *sql.DB) *PostgresTokenRepository {
	return &PostgresTokenRepository{DB: db}
}

func (this *PostgresTokenRepository) InsertRefreshToken(userID, clientIP, refreshToken string) error {
	_, err := this.DB.Exec(
		`INSERT INTO refresh_tokens (user_id, refresh_token_hash, client_ip)
		VALUES ($1, crypt($2, gen_salt('bf')), $3)`, userID, refreshToken, clientIP)

	return err
}

func (this *PostgresTokenRepository) CheckRefreshToken(userID, refreshToken string) (bool, error) {
	row := this.DB.QueryRow(
		`SELECT id FROM refresh_tokens
		WHERE user_id = $1
		AND refresh_token_hash = crypt($2, refresh_token_hash)
		AND is_used = false`)

	var id string
	err := row.Scan(&id)
	if err != nil {
		return false, nil
	}

	return err == nil, err
}

func (this *PostgresTokenRepository) MarkRefreshTokenUsed(userID, refreshToken string) error {
	_, err := this.DB.Exec(
		`UPDATE refresh_tokens
		SET is_used = true
		WHERE user_id = $1
		AND refresh_token_hash = crypt($2, refresh_token_hash)
		AND is_used = false`)
	return err
}
