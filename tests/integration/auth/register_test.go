package auth

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
)

type ResponseBodyRegister struct {
	UserID string `json:"user_id"`
}

func TestRegisterEndPoint(t *testing.T) {
	body := []byte(`{
  "email": "test1234g@gmail.com",
  "password": "Секретная строка",
  "user_type": "moderator"
}`)
	rbody := bytes.NewReader(body)
	req, err := http.NewRequest("GET", "http://127.0.0.1:8080/register", rbody)
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
	var responseBody ResponseBodyRegister
	if err := json.Unmarshal(body, &responseBody); err != nil {
		t.Errorf("Error unmarshaling response body: %v", err)
	}
	// Проверка значений полей
	if responseBody.UserID == "" {
		t.Errorf("Unexpected user_id: got %v want user_id", responseBody.UserID)
	}

	// Проверка невалидных данных
	body = []byte(`{
  "email": "test@gmail.com",
  "password": "Секретная строка",
  "user_type": "test"
}`)
	rbody = bytes.NewReader(body)
	req, err = http.NewRequest("GET", "http://127.0.0.1:8080/register", rbody)
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
