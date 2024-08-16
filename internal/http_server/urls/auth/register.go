package auth

import (
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"log/slog"
	"main/internal/storage/api"
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
			render.JSON(w, r, api.Error("failed to decode request."))
			return
		}
		log.Info("request body decoded", slog.Any("req", req))

		err = validator.New().Struct(req)
		if err != nil {
			validatorErr := err.(validator.ValidationErrors)
			log.Error("fail to validate body", "err", validatorErr)
			render.JSON(w, r, api.ValidationError(validatorErr))
			return
		}

		id, err := userRegister.CreateUser(req.Email, req.Password, req.UserType)
		if err != nil {
			log.Error("fail to create user", "err", err)
			render.JSON(w, r, "fail to create user")
			return
		}
		render.JSON(w, r, ResponseRegister{UserId: id})

	}
}
