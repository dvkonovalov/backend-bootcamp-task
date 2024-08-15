package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(storagePath string) (*Storage, error) {
	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s %s: %s", "Failed to open PostgreSQL storage at path", storagePath, err)
	}

	sqlRequest := `CREATE TABLE IF NOT EXISTS House(
		id SERIAL PRIMARY KEY,
		address TEXT NOT NULL,
		year INT NOT NULL,
    	developer TEXT,
    	created_at TIMESTAMP DEFAULT NOW(),
    	update_at TIMESTAMP DEFAULT NOW());
    `

	err = CreateTable(sqlRequest, db)
	if err != nil {
		return nil, fmt.Errorf("%s : %s", "Failed to create House table", err)
	}

	sqlRequest = `CREATE TABLE IF NOT EXISTS Apartments (
		id SERIAL PRIMARY KEY,
		price INT NOT NULL,
		rooms INT NOT NULL,
		house_id INT NOT NULL,
		status TEXT NOT NULL,
		FOREIGN KEY (house_id) REFERENCES House(id) ON DELETE CASCADE
	);
    `

	err = CreateTable(sqlRequest, db)
	if err != nil {
		return nil, fmt.Errorf("%s : %s", "Failed to create Apartments table", err)
	}

	return &Storage{db: db}, nil

}

func CreateTable(sqlRequest string, db *sql.DB) error {
	stmt, err := db.Prepare(sqlRequest)
	if err != nil {
		return fmt.Errorf("%s", err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("%s", err)
	}

	return nil
}
