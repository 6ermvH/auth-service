package repository

type IRefreshTokenRepository interface {
	Insert(userID, clientIP, refreshToken, accessTokenHash string) error
	Check(userID, refreshToken, accessTokenHash string) (bool, error)
	MarkUsed(userID, refreshToken string) error
	CleanupBadTokens() error
}
