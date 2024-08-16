package middleware

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"main/internal/storage/api"
	"net/http"
	"strings"
)

func CheckJWTToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	var mySecretKey []byte
	if authHeader == "" {
		return "", nil
	}

	tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

	keyFunc := func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signature method: %v", t.Header["alg"])
		}
		return mySecretKey, nil
	}

	// Разбор токена
	claims := &api.Claims{}
	parsedToken, err := jwt.ParseWithClaims(tokenString, claims, keyFunc)
	if err != nil {
		return "", fmt.Errorf("Error parse token: %s", err)
	}

	if !parsedToken.Valid {
		return "", fmt.Errorf("Invalid token")
	}

	return claims.Status, nil
}
