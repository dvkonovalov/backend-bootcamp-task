package db

import (
	"fmt"
	"main/internal/storage/api"
)

func (storage *Storage) CreateFlat(house_id int, price int, rooms int) (api.Flat, error) {
	var new_flat api.Flat
	stmt, err := storage.db.Prepare("INSERT INTO Apartments (house_id, price, rooms, status) VALUES ($1, $2, $3, 'created') RETURNING id, house_id, price, rooms, status;")
	if err != nil {
		return new_flat, fmt.Errorf("Error preparing statement: %s", err)
	}
	err = stmt.QueryRow(house_id, price, rooms).Scan(&new_flat.Id, &new_flat.House_id, &new_flat.Price, &new_flat.Rooms, &new_flat.Status)

	if err != nil {
		return new_flat, fmt.Errorf("Error executing query: %s", err)
	}
	err = storage.UpdateHouse(house_id)
	if err != nil {
		return new_flat, fmt.Errorf("Error updating house: %s", err)
	}
	return new_flat, nil
}

func (storage *Storage) UpdateFlat(id int, status string) (api.Flat, error) {
	var update_flat api.Flat
	stmt, err := storage.db.Prepare("UPDATE Apartments SET status=$1 WHERE id=$2 RETURNING id, house_id, price, rooms, status;")
	if err != nil {
		return update_flat, fmt.Errorf("Error preparing statement: %s", err)
	}
	err = stmt.QueryRow(status, id).Scan(&update_flat.Id, &update_flat.House_id, &update_flat.Price, &update_flat.Rooms, &update_flat.Status)

	if err != nil {
		return update_flat, fmt.Errorf("Error executing query: %s", err)
	}

	return update_flat, nil
}
