package middleware

import (
	"canonicalAuditlog/internal/jwtutil"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"strings"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Println("Missing Authorization header")
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Println("Invalid Authorization header format")
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenPart := parts[1]
		claims := &jwtutil.Claims{}

		// Todo: Remove hard coded fall back in case environment variable fails
		if len(jwtutil.JwtKey) == 0 {
			jwtutil.JwtKey = []byte("secret_key")
		}

		token, err := jwt.ParseWithClaims(tokenPart, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtutil.JwtKey, nil
		})

		if err != nil || !token.Valid {
			log.Println("Invalid JWT token:", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
