package responses

import (
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"net/http"
)

type ErrorResponse struct {
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
	Code      int    `json:"code"`
}

func ServerError(w http.ResponseWriter, r *http.Request, message string, code int) error {
	requestID := uuid.New().String()

	w.Header().Set("Retry-After", "10")
	w.WriteHeader(http.StatusInternalServerError)
	render.JSON(w, r, ErrorResponse{
		Message:   message,
		RequestID: requestID,
		Code:      code,
	})
	return nil
}

// 1 - ошибка декодирования тела запроса
// 10 - ошибка создания jwt-токена
// 11 - ошибка создания квартиры
// 12 - ошибка создания дома
// 20 - ошибка создания пользователя
// 21 - ошибка проверки пользователя
