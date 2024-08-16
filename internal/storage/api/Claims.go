package api

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	Username string `json:"username"`
	Status   string `json:"status"`
	jwt.RegisteredClaims
}
