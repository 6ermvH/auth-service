package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

const (
	accessTokenTTL  = 30 * time.Minute
	refreshTokenTTL = 7 * 24 * time.Hour
)

type TokenManager struct {
	jwtSecretKey []byte
}

type ITokenManager interface {
	GenerateNew(userID, clientIP string) (string, string, error)
	ParseAccess(accessToken, key string) (string, error)
	HashAccessToken(accessToken string) string
}

func NewTokenManager(secretKey []byte) *TokenManager {
	return &TokenManager{jwtSecretKey: secretKey}
}

func (tm *TokenManager) GenerateNew(userID, clientIP string) (string, string, error) {
	accessToken, err := tm.generateAccessToken(userID, clientIP)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := tm.generateRefreshToken()
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (tm *TokenManager) ParseAccess(accessToken, key string) (string, error) {
	token, err := jwt.Parse(accessToken, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS512 {
			return nil, fmt.Errorf("bad signing method %v", t.Header["alg"])
		}
		return tm.jwtSecretKey, nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrInvalidToken
	}

	value, ok := claims[key].(string)
	if !ok {
		return "", fmt.Errorf("'%v' key is missing in claims", key)
	}
	return value, nil
}

func (tm *TokenManager) HashAccessToken(accessToken string) string {
	hash := sha256.Sum256([]byte(accessToken))
	return hex.EncodeToString(hash[:])
}

func (tm *TokenManager) generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	refreshToken := base64.StdEncoding.EncodeToString(bytes)
	return refreshToken, nil
}

func (tm *TokenManager) generateAccessToken(userID, clientIP string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"ip":      clientIP,
		"exp":     time.Now().Add(accessTokenTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	accessToken, err := token.SignedString(tm.jwtSecretKey)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
