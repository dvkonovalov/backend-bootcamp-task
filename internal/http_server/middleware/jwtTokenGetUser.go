package middleware

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"main/internal/storage/api"
	"net/http"
	"os"
	"strings"
)

func CheckGetUser(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	mySecretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	if authHeader == "" {
		return "", fmt.Errorf("authorization header is missing")
	}

	tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

	keyFunc := func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signature method: %v", t.Header["alg"])
		}
		return mySecretKey, nil
	}

	claims := &api.Claims{}
	parsedToken, err := jwt.ParseWithClaims(tokenString, claims, keyFunc)
	if err != nil {
		return "", fmt.Errorf("error parse token: %s", err)
	}

	if !parsedToken.Valid {
		return "", fmt.Errorf("invalid token")
	}

	return claims.Username, nil
}
