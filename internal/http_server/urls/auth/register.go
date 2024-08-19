package auth

import (
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"log/slog"
	"main/internal/storage/api"
	"main/internal/storage/api/responses"
	"net/http"
)

type RequestRegister struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	UserType string `json:"user_type" validate:"required"`
}

type ResponseRegister struct {
	UserId string `json:"user_id"`
}

type UserRegister interface {
	CreateUser(email string, password string, userType string) (string, error)
}

func CreateUser(log *slog.Logger, userRegister UserRegister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestRegister
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
			http.Error(w, "Not found required data", http.StatusBadRequest)
			return
		}

		if req.UserType != api.Created && req.UserType != api.Moderator {
			log.Info("invalid user type in register", "type", req.UserType)
			http.Error(w, "invalid user type", http.StatusBadRequest)
			return
		}

		id, err := userRegister.CreateUser(req.Email, req.Password, req.UserType)
		if err != nil {
			log.Error("fail to create user", "err", err)
			err := responses.ServerError(w, r, "fail to create user", 20)
			if err != nil {
				log.Error("fail to send server error code", "err", err)
			}
			return
		}
		render.JSON(w, r, ResponseRegister{UserId: id})

	}
}
