package database

type IRefreshTokenRepository interface {
	InsertRefreshToken(userID, refreshToken, clientIP string) error
	CheckRefreshToken(userID, refreshToken string) (bool, error)
	MarkRefreshTokenUsed(userID, refreshToken string) error
}

