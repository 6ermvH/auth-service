package token

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateNew(userID, clientIP string) (string, string, error) {
	accessToken, err := generateAccessToken(userID, clientIP)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := generateRefreshToken()
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func ParseAccessToken(acess_token, key string) (string, error) {
	_, err := isOkAcessToken(acess_token)
	if err != nil {
		return "", err
	}
	token, _ := jwt.Parse(acess_token, func (t *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})
	claims, _ := token.Claims.(jwt.MapClaims)	

	value, ok := claims[key].(string)
	if !ok {
		return "", fmt.Errorf("'%v' key is missing in claims")
	}
	return value, nil
}

func isOkAcessToken(acess_token string) (bool, error) {
	token, err := jwt.Parse(acess_token, func (t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS512 {
			return nil, fmt.Errorf("Bad signing method %v", t.Header["alg"])
		}
		return jwtSecretKey, nil
	})
	if err != nil {
		return false, err
	}

	if !token.Valid {
		return false, fmt.Errorf("invalid token signature")
	}

	return true, nil
}
