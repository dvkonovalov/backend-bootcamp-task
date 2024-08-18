package auth

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"main/internal/config"
	"net/http"
	"testing"
)

type ResponseBody struct {
	Token string `json:"token"`
}

func TestDummyLoginEndPoint(t *testing.T) {
	cnf := config.MustLoad()
	body := []byte(`{"user_type":"client"}`)
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
	// Чтение тела ответа
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: %v", err)
	}
	// Декодирование JSON-ответа
	var responseBody ResponseBody
	if err := json.Unmarshal(body, &responseBody); err != nil {
		t.Errorf("Error unmarshaling response body: %v", err)
	}
	// Проверка значений полей
	if responseBody.Token == "" {
		t.Errorf("Unexpected token: got %v want token", responseBody.Token)
	}

	// ПРОВЕРКА МОДЕРАТОРА
	body = []byte(`{"user_type":"moderator"}`)
	rbody = bytes.NewReader(body)
	req, err = http.NewRequest("GET", "http://"+cnf.HttpServer.Address+"/dummyLogin", rbody)
	if err != nil {
		t.Errorf("Error creating http request: %v", err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("Error making http request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Unexpected http status code: %v", resp.StatusCode)
	}
	// Чтение тела ответа
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: %v", err)
	}
	// Декодирование JSON-ответа
	if err := json.Unmarshal(body, &responseBody); err != nil {
		t.Errorf("Error unmarshaling response body: %v", err)
	}
	// Проверка значений полей
	if responseBody.Token == "" {
		t.Errorf("Unexpected token: got %v want token", responseBody.Token)
	}

	// НЕВЕРНЫЕ ДАННЫЕ
	body = []byte(`{"user_type":"test"}`)
	rbody = bytes.NewReader(body)
	req, err = http.NewRequest("GET", "http://"+cnf.HttpServer.Address+"/dummyLogin", rbody)
	if err != nil {
		t.Errorf("Error creating http request: %v", err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("Error making http request: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Unexpected http status code: %v", resp.StatusCode)
	}

}
