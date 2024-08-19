package auth

import (
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"log/slog"
	"main/internal/http_server/middleware"
	"main/internal/storage/api"
	"main/internal/storage/api/responses"
	"net/http"
)

type RequestDummyLogin struct {
	UserType string `json:"user_type" validate:"required"`
}

type ResponseDummyLogin struct {
	Token string `json:"token"`
}

func CreateToken(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestDummyLogin
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("fail to decode body", "err", err)
			err := responses.ServerError(w, r, "fail to decode body", 1)
			if err != nil {
				log.Error("fail to send server error code", "err", err)
			}
			return
		}
		log.Info("request body decoded", slog.Any("req", req))

		err = validator.New().Struct(req)
		if err != nil {
			var validatorErr validator.ValidationErrors
			errors.As(err, &validatorErr)
			log.Error("fail to validate body", "err", validatorErr)
			http.Error(w, "Not found user type", http.StatusBadRequest)
			return
		}

		if req.UserType != api.Client && req.UserType != api.Moderator {
			log.Warn("Invalid user type", "type", req.UserType)
			http.Error(w, "Invalid user type", http.StatusBadRequest)
			return
		}

		tokenString, err := middleware.CreateJWTToken("Simple", req.UserType)
		if err != nil {
			log.Error("fail to create token", "err", err)
			err := responses.ServerError(w, r, "fail to create token", 10)
			if err != nil {
				log.Error("fail to send server error code", "err", err)
			}
			return
		}

		render.JSON(w, r, ResponseDummyLogin{Token: tokenString})
	}
}
