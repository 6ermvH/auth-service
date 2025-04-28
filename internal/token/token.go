package token

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenManager struct {
	jwtSecretKey []byte
}

type ITokenManager interface {
	GenerateNew(userID, clientIP string) (string, string, error)
	ParseAccess(accessToken, key string) (string, error)
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

func (tm *TokenManager) ParseAccess(access_token, key string) (string, error) {
	_, err := tm.isOkAccess(access_token)
	if err != nil {
		return "", err
	}
	token, _ := jwt.Parse(access_token, func(t *jwt.Token) (interface{}, error) {
		return tm.jwtSecretKey, nil
	})
	claims, _ := token.Claims.(jwt.MapClaims)

	value, ok := claims[key].(string)
	if !ok {
		return "", fmt.Errorf("'%v' key is missing in claims", key)
	}
	return value, nil
}

func (tm *TokenManager) isOkAccess(access_token string) (bool, error) {
	token, err := jwt.Parse(access_token, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS512 {
			return nil, fmt.Errorf("Bad signing method %v", t.Header["alg"])
		}
		return tm.jwtSecretKey, nil
	})
	if err != nil {
		return false, err
	}

	if !token.Valid {
		return false, fmt.Errorf("invalid token signature")
	}

	return true, nil
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
		"exp":     time.Now().Add(30 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	accessToken, err := token.SignedString(tm.jwtSecretKey)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
