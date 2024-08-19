package auth

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestLoginEndPoint(t *testing.T) {
	password := "Секретная строка"
	body := []byte(`{
  "email": "test32@gmail.com",
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

	body = []byte(`{
	  "id": "` + responseBody.UserID + `",
	  "password": "` + password + `"
	}`)
	rbody = bytes.NewReader(body)
	req, err = http.NewRequest("GET", "http://127.0.0.1:8080/login", rbody)
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
	var responseBodyToken ResponseBody
	if err := json.Unmarshal(body, &responseBodyToken); err != nil {
		t.Errorf("Error unmarshaling response body: %v", err)
	}
	// Проверка значений полей
	if responseBodyToken.Token == "" {
		t.Errorf("Unexpected user_id: got %v want token", responseBodyToken.Token)
	}

	// Несуществующий пользователь
	body = []byte(`{
	  "id": "1000000000000",
	  "password": "` + password + `"
	}`)
	rbody = bytes.NewReader(body)
	req, err = http.NewRequest("GET", "http://127.0.0.1:8080/login", rbody)
	if err != nil {
		t.Errorf("Error creating http request: %v", err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("Error making http request: %v", err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Unexpected http status code: %v", resp.StatusCode)
	}

}
