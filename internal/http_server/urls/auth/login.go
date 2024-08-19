package auth

import (
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"log/slog"
	"net/http"
)

type RequestLogin struct {
	Id       string `json:"id" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ResponseLogin struct {
	Token string `json:"token"`
}

type UserLogin interface {
	LoginUser(id string, password string) (string, error)
}

func LoginUser(log *slog.Logger, userLogin UserLogin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestLogin
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("fail to decode body", "err", err)
			http.Error(w, "fail to decode body", http.StatusInternalServerError)
			return
		}
		log.Info("request body decoded", slog.Any("req", req))

		err = validator.New().Struct(req)
		if err != nil {
			var validatorErr validator.ValidationErrors
			errors.As(err, &validatorErr)
			log.Error("fail to validate body", "err", validatorErr)
			http.Error(w, "fail to validate body", http.StatusBadRequest)
			return
		}

		jwtToken, err := userLogin.LoginUser(req.Id, req.Password)
		if err != nil {
			log.Error("fail to login", "err", err)
			http.Error(w, "fail to login", http.StatusNotFound)
			return
		}
		render.JSON(w, r, ResponseLogin{Token: jwtToken})

	}

}
