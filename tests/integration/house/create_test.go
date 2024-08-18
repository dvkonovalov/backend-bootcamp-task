package house

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
)

type ResponseBody struct {
	Token string `json:"token"`
}

type ResponseBodyCreate struct {
	Id        int    `json:"id"`
	Address   string `json:"address"`
	Year      int    `json:"year"`
	Developer string `json:"developer"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func TestCreateEndPoint(t *testing.T) {
	body := []byte(`{"user_type":"moderator"}`)
	rbody := bytes.NewReader(body)
	req, err := http.NewRequest("GET", "http://0.0.0.0:8080/dummyLogin", rbody)
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

	// Нормальный запрос
	body = []byte(`{
	  "address": "Лесная улица, 7, Москва, 125196",
	  "year": 2000,
	  "developer": "Мэрия города"
	}`)
	rbody = bytes.NewReader(body)
	req, err = http.NewRequest("GET", "http://0.0.0.0:8080/house/create", rbody)
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
	if responseBodyCreate.CreatedAt == "" {
		t.Errorf("Unexpected id: got %v want id", responseBodyCreate.Id)
	}
	if responseBodyCreate.Address != "Лесная улица, 7, Москва, 125196" {
		t.Errorf("Unexpected address: got %v want address", responseBodyCreate.Id)
	}
	if responseBodyCreate.Year < 0 {
		t.Errorf("Unexpected Year: got %v want Year", responseBodyCreate.Id)
	}
	if responseBodyCreate.Developer != "Мэрия города" {
		t.Errorf("Unexpected Developer: got %v want Developer", responseBodyCreate.Id)
	}

	// Неавторизованный доступ
	req, err = http.NewRequest("GET", "http://0.0.0.0:8080/house/create", rbody)
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
