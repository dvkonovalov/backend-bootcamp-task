package db

import (
	"fmt"
	"main/internal/storage/api"
)

func (storage *Storage) CreateFlat(houseId int, price int, rooms int) (api.Flat, error) {
	var newFlat api.Flat
	stmt, err := storage.db.Prepare("INSERT INTO Apartments (house_id, price, rooms, status) VALUES ($1, $2, $3, 'created') RETURNING id, house_id, price, rooms, status;")
	if err != nil {
		return newFlat, fmt.Errorf("error preparing statement: %s", err)
	}
	err = stmt.QueryRow(houseId, price, rooms).Scan(&newFlat.Id, &newFlat.HouseId, &newFlat.Price, &newFlat.Rooms, &newFlat.Status)

	if err != nil {
		return newFlat, fmt.Errorf("error executing query: %s", err)
	}
	err = storage.UpdateHouse(houseId)
	if err != nil {
		return newFlat, fmt.Errorf("error updating house: %s", err)
	}
	return newFlat, nil
}

func (storage *Storage) UpdateFlat(id int, status string, moderator string) (api.Flat, error) {
	var updateFlat api.Flat

	statusNow, err := storage.GetStatus(id)
	if err != nil {
		return updateFlat, fmt.Errorf("error getting status for id %d: %s", id, err)
	}

	if statusNow == api.OnModeration {
		var moderatorNow string
		stmt, err := storage.db.Prepare("SELECT moderator FROM Moderation WHERE flat_id=$1;")
		if err != nil {
			return updateFlat, fmt.Errorf("error preparing statement: %s", err)
		}
		err = stmt.QueryRow(id).Scan(&moderatorNow)
		if err != nil {
			return updateFlat, fmt.Errorf("error executing query: %s", err)
		}
		if moderator != moderatorNow {
			return updateFlat, fmt.Errorf("this apartment is being moderated by %s, but %s wants to moderate it", moderatorNow, moderator)
		}

		stmt, err = storage.db.Prepare("UPDATE Apartments SET status=$1 WHERE id=$2 RETURNING id, house_id, price, rooms, status;")
		if err != nil {
			return updateFlat, fmt.Errorf("error preparing statement: %s", err)
		}
		err = stmt.QueryRow(status, id).Scan(&updateFlat.Id, &updateFlat.HouseId, &updateFlat.Price, &updateFlat.Rooms, &updateFlat.Status)

		if err != nil {
			return updateFlat, fmt.Errorf("error executing query: %s", err)
		}

		return updateFlat, nil

	} else {
		stmt, err := storage.db.Prepare("UPDATE Apartments SET status=$1 WHERE id=$2 RETURNING id, house_id, price, rooms, status;")
		if err != nil {
			return updateFlat, fmt.Errorf("error preparing statement: %s", err)
		}
		err = stmt.QueryRow(status, id).Scan(&updateFlat.Id, &updateFlat.HouseId, &updateFlat.Price, &updateFlat.Rooms, &updateFlat.Status)

		if err != nil {
			return updateFlat, fmt.Errorf("error executing query: %s", err)
		}

		return updateFlat, nil
	}

}

func (storage *Storage) GetStatus(id int) (string, error) {
	var status string
	stmt, err := storage.db.Prepare("SELECT status FROM Apartments WHERE id=$1;")
	if err != nil {
		return "", fmt.Errorf("error preparing statement: %s", err)
	}
	err = stmt.QueryRow(id).Scan(&status)

	if err != nil {
		return "", fmt.Errorf("error executing query: %s", err)
	}

	return status, nil

}

func (storage *Storage) BlockModerationOtherAdmin(flatId int, moderator string) (bool, error) {
	stmt, err := storage.db.Prepare("INSERT INTO Moderation(flat_id, moderator) VALUES ($1, $2);")
	if err != nil {
		return false, fmt.Errorf("error preparing statement: %s", err)
	}
	_, err = stmt.Exec(flatId, moderator)

	if err != nil {
		return false, fmt.Errorf("error executing query: %s", err)
	}

	return true, nil

}
