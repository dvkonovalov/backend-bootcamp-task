package db

import (
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"main/internal/http_server/middleware"
)

func (storage *Storage) CreateUser(email string, password string, userType string) (string, error) {
	var id string
	var count = 0
	stmt, err := storage.db.Prepare("SELECT id from Users WHERE email = $1;")
	if err != nil {
		return id, fmt.Errorf("Error preparing statement: %s", err)
	}
	err = stmt.QueryRow(email).Scan(&id)

	if err != nil && err != sql.ErrNoRows {
		return id, fmt.Errorf("Error executing query: %s", err)
	}
	if count > 0 {
		return id, fmt.Errorf("User already exists")
	}

	stmt, err = storage.db.Prepare("INSERT INTO Users (email, password_hash, user_type) VALUES ($1, $2, $3) RETURNING id;")
	if err != nil {
		return id, fmt.Errorf("Error preparing statement: %s", err)
	}
	if password == "" {
		return id, fmt.Errorf("Password is empty")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return id, fmt.Errorf("Error in calc Hash password", "err", err)
	}
	err = stmt.QueryRow(email, hashedPassword, userType).Scan(&id)

	if err != nil {
		return id, fmt.Errorf("Error executing query: %s", err)
	}

	return id, nil
}

func (storage *Storage) LoginUser(id string, password string) (string, error) {
	var find_id int
	var password_hash, userType string
	stmt, err := storage.db.Prepare("SELECT id, password_hash, user_type from Users WHERE id = $1;")
	if err != nil {
		return "", fmt.Errorf("Error preparing statement: %s", err)
	}
	err = stmt.QueryRow(id).Scan(&find_id, &password_hash, &userType)

	if err != nil {
		return "", fmt.Errorf("Error executing query: %s", err)
	}
	if find_id == 0 {
		return "", fmt.Errorf("User not found")
	}

	if password == "" {
		return "", fmt.Errorf("Password is empty")
	}
	err = bcrypt.CompareHashAndPassword([]byte(password_hash), []byte(password))
	if err != nil {
		return "", fmt.Errorf("Invalid password", "err", err)
	}
	jwtToken, err := middleware.CreateJWTToken(id, userType)
	if err != nil {
		return "", fmt.Errorf("Error generating JWT token", "err", err)
	}
	return jwtToken, nil
}
