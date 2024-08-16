package auth

import (
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"log/slog"
	"main/internal/http_server/middleware"
	"main/internal/storage/api"
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

		if req.UserType != "moderator" && req.UserType != "client" {
			log.Warn("Invalid user type", "type", req.UserType)
			render.JSON(w, r, "Invalid user type.")
			return
		}

		tokenString, err := middleware.CreateJWTToken("Simple", req.UserType)
		if err != nil {
			log.Error("fail to create token", "err", err)
			render.JSON(w, r, api.Error("failed to create token"))
			return
		}

		render.JSON(w, r, ResponseDummyLogin{Token: tokenString})
	}
}
