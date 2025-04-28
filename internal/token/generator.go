package token

import (
	"crypto/rand"
	"encoding/base64"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecretKey = []byte(os.Getenv("JWT_SECRET_KEY"))
)

func generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	refreshToken := base64.StdEncoding.EncodeToString(bytes)
	return refreshToken, nil
}

func generateAccessToken(userID, clientIP string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"ip":      clientIP,
		"exp":     time.Now().Add(30 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	accessToken, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
