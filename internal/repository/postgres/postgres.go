package postgres

import (
	"database/sql"

	_ "example.com/auth_service/repository"
	_ "github.com/lib/pq"
)

type PostgresTokenRepository struct {
	DB *sql.DB
}

func NewRefreshTokenRepository(db *sql.DB) *PostgresTokenRepository {
	return &PostgresTokenRepository{DB: db}
}

func (this *PostgresTokenRepository) Insert(userID, clientIP, refreshToken, accessTokenHash string) error {
	_, err := this.DB.Exec(
		`INSERT INTO refresh_tokens (user_id, client_ip, refresh_token_hash, access_token_sha256)
		VALUES ($1, $2, crypt($3, gen_salt('bf')), $4)`, userID, clientIP, refreshToken, accessTokenHash)

	return err
}

func (this *PostgresTokenRepository) Check(userID, refreshToken, accessTokenHash string) (bool, error) {
	this.CleanupBadTokens()
	row := this.DB.QueryRow(
		`SELECT id FROM refresh_tokens
		WHERE user_id = $1
		AND refresh_token_hash = crypt($2, refresh_token_hash)
		AND access_token_sha256 = $3
		AND is_used = false`, userID, refreshToken, accessTokenHash)

	var id string
	err := row.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return err == nil, err
}

func (this *PostgresTokenRepository) MarkUsed(userID, refreshToken string) error {
	_, err := this.DB.Exec(
		`UPDATE refresh_tokens
		SET is_used = true
		WHERE user_id = $1
		AND refresh_token_hash = crypt($2, refresh_token_hash)
		AND is_used = false`, userID, refreshToken)
	return err
}

func (this *PostgresTokenRepository) CleanupBadTokens() error {
	_, err := this.DB.Exec(
		`DELETE FROM refresh_tokens
		WHERE is_used = true
		OR expired_at <= NOW()`)
	return err
}
