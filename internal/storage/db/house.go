package db

import (
	"context"
	"fmt"
	"main/internal/storage/api"
)

func (storage *Storage) CreateHouse(address string, developer string, year int) (api.House, error) {
	var newHouse api.House
	stmt, err := storage.Db.Prepare("INSERT INTO House(address, developer, year, created_at, update_at) VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id, address, developer, year, created_at, update_at ")
	if err != nil {
		return newHouse, fmt.Errorf("error in CreateHouse Prepare: %s", err)
	}

	err = stmt.QueryRow(address, developer, year).Scan(&newHouse.Id, &newHouse.Address, &newHouse.Developer, &newHouse.Year, &newHouse.DateCreate, &newHouse.UpdateAt)
	if err != nil {
		return newHouse, fmt.Errorf("error in CreateHouse Exec: %s", err)
	}

	return newHouse, nil
}

func (storage *Storage) UpdateHouse(houseId int) error {
	stmt, err := storage.Db.Prepare("UPDATE House SET update_at = Now() WHERE id=$1;")
	if err != nil {
		return fmt.Errorf("error preparing statement: %s", err)
	}
	_, err = stmt.Exec(houseId)
	if err != nil {
		return fmt.Errorf("error executing query: %s", err)
	}
	return nil
}

func (storage *Storage) GetAllFlats(houseId int, userType string) ([]api.Flat, error) {
	var flats []api.Flat
	switch userType {
	case api.Moderator:
		{
			stmt, err := storage.Db.Prepare("SELECT id, house_id, price, rooms, status FROM Apartments WHERE house_id=$1")
			if err != nil {
				return flats, fmt.Errorf("error in Prapare request in GetAllFlats: %s", err)
			}
			rows, err := stmt.Query(houseId)
			if err != nil {
				return flats, fmt.Errorf("error Exec in GetAllFlats: %s", err)
			}
			for rows.Next() {
				var flat api.Flat
				err = rows.Scan(&flat.Id, &flat.HouseId, &flat.Price, &flat.Rooms, &flat.Status)
				if err != nil {
					return flats, fmt.Errorf("error getting data from a row in GetAllFlats: %s", err)
				}
				flats = append(flats, flat)
			}
			return flats, nil
		}
	case api.Client:
		{
			stmt, err := storage.Db.Prepare("SELECT id, house_id, price, rooms, status FROM Apartments WHERE house_id=$1 AND status=$2")
			if err != nil {
				return flats, fmt.Errorf("error in Prapare request in GetAllFlats: %s", err)
			}
			rows, err := stmt.Query(houseId, api.Approved)
			if err != nil {
				return flats, fmt.Errorf("error Exec in GetAllFlats: %s", err)
			}
			for rows.Next() {
				var flat api.Flat
				err = rows.Scan(&flat.Id, &flat.HouseId, &flat.Price, &flat.Rooms, &flat.Status)
				if err != nil {
					return flats, fmt.Errorf("error getting data from a row in GetAllFlats: %s", err)
				}
				flats = append(flats, flat)
			}
			return flats, nil
		}
	default:
		{
			stmt, err := storage.Db.Prepare("SELECT id, house_id, price, rooms, status FROM Apartments WHERE house_id=$1 AND status=$2")
			if err != nil {
				return flats, fmt.Errorf("error in Prapare request in GetAllFlats: %s", err)
			}
			rows, err := stmt.Query(houseId, api.Approved)
			if err != nil {
				return flats, fmt.Errorf("error Exec in GetAllFlats: %s", err)
			}
			for rows.Next() {
				var flat api.Flat
				err = rows.Scan(&flat.Id, &flat.HouseId, &flat.Price, &flat.Rooms, &flat.Status)
				if err != nil {
					return flats, fmt.Errorf("error getting data from a row in GetAllFlats: %s", err)
				}
				flats = append(flats, flat)
			}
			return flats, nil

		}

	}
}

func (storage *Storage) SubscribeToHouse(ctx context.Context, houseID string, email string) error {
	_, err := storage.Db.ExecContext(ctx, `
		INSERT INTO house_subscriptions (house_id, user_email)
		VALUES ($1, $2)
	`, houseID, email)
	return err
}

func (storage *Storage) GetHouseSubscribers(ctx context.Context, houseID string) ([]string, error) {
	rows, err := storage.Db.QueryContext(ctx, `
		SELECT user_email FROM house_subscriptions WHERE house_id = $1
	`, houseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var emails []string
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			return nil, err
		}
		emails = append(emails, email)
	}

	return emails, nil
}
