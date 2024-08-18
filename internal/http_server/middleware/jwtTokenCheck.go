package middleware

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"main/internal/storage/api"
	"net/http"
	"os"
	"strings"
)

func CheckJWTToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	//mySecretKey := []byte("secret")
	mySecretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	if authHeader == "" {
		return "", fmt.Errorf("Authorization header is missing")
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
