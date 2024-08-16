package auth

import (
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"log/slog"
	"main/internal/storage/api"
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

		jwtToken, err := userLogin.LoginUser(req.Id, req.Password)
		if err != nil {
			log.Error("fail to login", "err", err)
			render.JSON(w, r, "failed to login.")
			return
		}
		render.JSON(w, r, ResponseLogin{Token: jwtToken})

	}

}
