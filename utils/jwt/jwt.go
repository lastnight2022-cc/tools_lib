package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateJWTToken(userId uint, secret string, expireTime time.Duration) (string, error) {
	if userId == 0 {
		return "", errors.New("invalid user ID")
	}
	if secret == "" {
		return "", errors.New("secret cannot be empty")
	}
	if expireTime <= 0 {
		return "", errors.New("expire time must be positive")
	}
	expirationTime := time.Now().Add(expireTime)
	claims := &jwt.StandardClaims{
		Subject:   fmt.Sprintf("%d", userId),
		ExpiresAt: expirationTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func VerifyJWTToken(tokenString string, secret string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%s", "Unexpected signing method")
		}
		return []byte(secret), nil
	})
}
