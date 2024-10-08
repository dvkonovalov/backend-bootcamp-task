package middleware

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"main/internal/storage/api"
	"os"
)

func CreateJWTToken(username string, userType string) (string, error) {
	mySecretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	//mySecretKey := []byte("secret")
	claims := api.Claims{Status: userType, Username: username}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен нашим секретным ключем
	tokenString, err := token.SignedString(mySecretKey)
	if err != nil {
		return "", fmt.Errorf("fail to sign token. Error: %s", err)
	}
	return tokenString, nil
}
