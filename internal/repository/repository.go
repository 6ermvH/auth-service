package repository

type IRefreshTokenRepository interface {
	Insert(userID, refreshToken, clientIP string) error
	Check(userID, refreshToken string) (bool, error)
	MarkUsed(userID, refreshToken string) error
	CleanupBadTokens() error
}
