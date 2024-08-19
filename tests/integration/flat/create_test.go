package flat

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"main/internal/storage/api"
	"net/http"
	"testing"
)

type ResponseBodyCreate struct {
	Id      int    `json:"id"`
	HouseId int    `json:"house_id"`
	Price   int    `json:"price"`
	Rooms   int    `json:"rooms"`
	Status  string `json:"status"`
}

func TestCreateEndPoint(t *testing.T) {
	body := []byte(`{"user_type":"moderator"}`)
	rbody := bytes.NewReader(body)
	req, err := http.NewRequest("GET", "http://127.0.0.1:8080/dummyLogin", rbody)
	if err != nil {
		t.Errorf("Error creating http request: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("Error making http request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Unexpected http status code: %v", resp.StatusCode)
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: %v", err)
	}
	var responseBody ResponseBody
	if err := json.Unmarshal(body, &responseBody); err != nil {
		t.Errorf("Error unmarshaling response body: %v", err)
	}
	if responseBody.Token == "" {
		t.Errorf("Unexpected token: got %v want token", responseBody.Token)
	}

	// Создаем дом
	body = []byte(`{
	  "address": "Лесная улица, 7, Москва, 125196",
	  "year": 2000,
	  "developer": "Мэрия города"
	}`)
	rbody = bytes.NewReader(body)
	req, err = http.NewRequest("GET", "http://127.0.0.1:8080/house/create", rbody)
	if err != nil {
		t.Errorf("Error creating http request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+responseBody.Token)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("Error making http request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Unexpected http status code: %v", resp.StatusCode)
	}

	//Нормальный запрос
	body = []byte(`{
	  "house_id": 1,
	  "price": 10000,
	  "rooms": 4
	}`)
	rbody = bytes.NewReader(body)
	req, err = http.NewRequest("GET", "http://127.0.0.1:8080/flat/create", rbody)
	if err != nil {
		t.Errorf("Error creating http request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+responseBody.Token)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("Error making http request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Unexpected http status code: %v", resp.StatusCode)
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: %v", err)
	}
	var responseBodyCreate ResponseBodyCreate
	if err := json.Unmarshal(body, &responseBodyCreate); err != nil {
		t.Errorf("Error unmarshaling response body: %v", err)
	}
	if responseBodyCreate.Id < 0 {
		t.Errorf("Unexpected id: got %v want id", responseBodyCreate.Id)
	}
	if responseBodyCreate.Price < 0 {
		t.Errorf("Unexpected Price: got %v want Price", responseBodyCreate.Price)
	}
	if responseBodyCreate.Rooms < 0 {
		t.Errorf("Unexpected Rooms: got %v want Rooms", responseBodyCreate.Rooms)
	}
	if responseBodyCreate.HouseId < 0 {
		t.Errorf("Unexpected HouseId: got %v want HouseId", responseBodyCreate.HouseId)
	}
	if responseBodyCreate.Status != api.Created {
		t.Errorf("Unexpected Status: got %v want Status", responseBodyCreate.Status)
	}

	// Неавторизованный доступ
	req, err = http.NewRequest("GET", "http://127.0.0.1:8080/flat/create", rbody)
	if err != nil {
		t.Errorf("Error creating http request: %v", err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("Error making http request: %v", err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Unexpected http status code: %v", resp.StatusCode)
	}

}
