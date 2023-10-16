package jwtutil

import (
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

// JwtKey Todo: An environment variable but a hard coded fall back has been added
var JwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

type Claims struct {
	jwt.StandardClaims
}

func GenerateToken() (string, error) {
	expirationTime := time.Now().Add(336 * time.Hour)
	claims := &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Todo: Remove hard coded fall back in case environment variable fails
	if len(JwtKey) == 0 {
		JwtKey = []byte("secret_key")
	}

	return token.SignedString(JwtKey)
}
