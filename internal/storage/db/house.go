package db

import (
	"fmt"
	"main/internal/storage/api"
)

func (storage *Storage) CreateHouse(address string, developer string, year int) (api.House, error) {
	var new_house api.House
	stmt, err := storage.db.Prepare("INSERT INTO House(address, developer, year, created_at, update_at) VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id, address, developer, year, created_at, update_at ")
	if err != nil {
		return new_house, fmt.Errorf("Error in CreateHouse Prepare: %s", err)
	}

	err = stmt.QueryRow(address, developer, year).Scan(&new_house.Id, &new_house.Address, &new_house.Developer, &new_house.Year, &new_house.CreatedAt, &new_house.UpdateAt)
	if err != nil {
		return new_house, fmt.Errorf("Error in CreateHouse Exec: %s", err)
	}

	return new_house, nil
}

func (storage *Storage) UpdateHouse(house_id int) error {
	stmt, err := storage.db.Prepare("UPDATE House SET update_at = Now() WHERE id=$1;")
	if err != nil {
		return fmt.Errorf("Error preparing statement: %s", err)
	}
	_, err = stmt.Exec(house_id)
	if err != nil {
		return fmt.Errorf("Error executing query: %s", err)
	}
	return nil
}

func (storage *Storage) GetAllFlats(house_id int) ([]api.Flat, error) {
	var flats []api.Flat
	stmt, err := storage.db.Prepare("SELECT id, house_id, price, rooms, status FROM Apartments WHERE house_id=$1")
	if err != nil {
		return flats, fmt.Errorf("Error in Prapare request in GetAllFlats: %s", err)
	}
	rows, err := stmt.Query(house_id)
	if err != nil {
		return flats, fmt.Errorf("Error Exec in GetAllFlats: %s", err)
	}
	for rows.Next() {
		var flat api.Flat
		err = rows.Scan(&flat.Id, &flat.House_id, &flat.Price, &flat.Rooms, &flat.Status)
		if err != nil {
			return flats, fmt.Errorf("Error getting data from a row in GetAllFlats: %s", err)
		}
		flats = append(flats, flat)
	}
	return flats, nil
}
