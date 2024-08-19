package db

import (
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"main/internal/http_server/middleware"
)

func (storage *Storage) CreateUser(email string, password string, userType string) (string, error) {
	var id string
	stmt, err := storage.Db.Prepare("SELECT id from Users WHERE email = $1;")
	if err != nil {
		return id, fmt.Errorf("error preparing statement: %s", err)
	}
	err = stmt.QueryRow(email).Scan(&id)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return id, fmt.Errorf("error executing query: %s", err)
	}
	if id != "" {
		return id, fmt.Errorf("user already exists")
	}

	stmt, err = storage.Db.Prepare("INSERT INTO Users (email, password_hash, user_type) VALUES ($1, $2, $3) RETURNING id;")
	if err != nil {
		return id, fmt.Errorf("error preparing statement: %s", err)
	}
	if password == "" {
		return id, fmt.Errorf("password is empty")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return id, fmt.Errorf("error in calc Hash password. Error: %s", err)
	}
	err = stmt.QueryRow(email, hashedPassword, userType).Scan(&id)

	if err != nil {
		return id, fmt.Errorf("error executing query: %s", err)
	}

	return id, nil
}

func (storage *Storage) LoginUser(id string, password string) (string, error) {
	var findId int
	var passwordHash, userType string
	stmt, err := storage.Db.Prepare("SELECT id, password_hash, user_type from Users WHERE id = $1;")
	if err != nil {
		return "", fmt.Errorf("error preparing statement: %s", err)
	}
	err = stmt.QueryRow(id).Scan(&findId, &passwordHash, &userType)

	if err != nil {
		return "", fmt.Errorf("error executing query: %s", err)
	}
	if findId == 0 {
		return "", fmt.Errorf("user not found")
	}

	if password == "" {
		return "", fmt.Errorf("password is empty")
	}
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return "", fmt.Errorf("invalid password. Error: %s", err)
	}
	jwtToken, err := middleware.CreateJWTToken(id, userType)
	if err != nil {
		return "", fmt.Errorf("error generating JWT token. Error: %s", err)
	}
	return jwtToken, nil
}
