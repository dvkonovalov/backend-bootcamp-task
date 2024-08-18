package house

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"main/internal/config"
	"net/http"
	"strconv"
	"testing"
)

type ResponseBody struct {
	Token string `json:"token"`
}

type ResponseBodyCreateHouse struct {
	Id        int    `json:"id"`
	Address   string `json:"address"`
	Year      int    `json:"year"`
	Developer string `json:"developer"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type ResponseBodyCreateFlat struct {
	Id      int    `json:"id"`
	HouseId int    `json:"house_id"`
	Price   int    `json:"price"`
	Rooms   int    `json:"rooms"`
	Status  string `json:"status"`
}

type ResponseBodyGetFlat struct {
	Flats []ResponseBodyCreateFlat `json:"flats"`
}

func TestGetFlatsEndPoint(t *testing.T) {
	cnf := config.MustLoad()
	// Получаем токен
	body := []byte(`{"user_type":"moderator"}`)
	rbody := bytes.NewReader(body)
	req, err := http.NewRequest("GET", "http://"+cnf.HttpServer.Address+"/dummyLogin", rbody)
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
	req, err = http.NewRequest("GET", "http://"+cnf.HttpServer.Address+"/house/create", rbody)
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
	var responseBodyCreateHouse ResponseBodyCreateHouse
	if err := json.Unmarshal(body, &responseBodyCreateHouse); err != nil {
		t.Errorf("Error unmarshaling response body: %v", err)
	}

	// Создаем квартиры в нем
	body = []byte(`{
	  "house_id": ` + strconv.Itoa(responseBodyCreateHouse.Id) + `,
	  "price": 10000,
	  "rooms": 4
	}`)
	rbody = bytes.NewReader(body)
	req, err = http.NewRequest("GET", "http://"+cnf.HttpServer.Address+"/flat/create", rbody)
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
	var responseBodyCreateFlat ResponseBodyCreateFlat
	if err := json.Unmarshal(body, &responseBodyCreateFlat); err != nil {
		t.Errorf("Error unmarshaling response body: %v", err)
	}
	body = []byte(`{
	  "house_id": ` + strconv.Itoa(responseBodyCreateHouse.Id) + `,
	  "price": 20000,
	  "rooms": 3
	}`)
	rbody = bytes.NewReader(body)
	req, err = http.NewRequest("GET", "http://"+cnf.HttpServer.Address+"/flat/create", rbody)
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
	var responseBodyCreateFlat2 ResponseBodyCreateFlat
	if err := json.Unmarshal(body, &responseBodyCreateFlat2); err != nil {
		t.Errorf("Error unmarshaling response body: %v", err)
	}

	// Нормальный запрос
	req, err = http.NewRequest("GET", "http://"+cnf.HttpServer.Address+"/house/"+strconv.Itoa(responseBodyCreateHouse.Id), nil)
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
	var responseBodyGetFlat ResponseBodyGetFlat
	if err := json.Unmarshal(body, &responseBodyGetFlat); err != nil {
		t.Errorf("Error unmarshaling response body: %v", err)
	}

	if len(responseBodyGetFlat.Flats) != 2 {
		t.Errorf("Unexpected count of flats: got %d, want 2", len(responseBodyGetFlat.Flats))
	}
	if responseBodyGetFlat.Flats[0].Id != responseBodyCreateFlat.Id {
		t.Errorf("Error id flats: got %v, want %v", responseBodyGetFlat.Flats[0].Id, responseBodyCreateFlat.Id)
	}
	if responseBodyGetFlat.Flats[1].Id != responseBodyCreateFlat2.Id {
		t.Errorf("Error id flats: got %v, want %v", responseBodyGetFlat.Flats[1].Id, responseBodyCreateFlat2.Id)
	}

	// Неавторизованный доступ
	req, err = http.NewRequest("GET", "http://"+cnf.HttpServer.Address+"/house/"+strconv.Itoa(responseBodyCreateHouse.Id), nil)
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
