package house

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

type Request struct {
	Address   string `json:"address" validate:"required"`
	Developer string `json:"developer,omitempty"`
	Year      int    `json:"year"  validate:"required"`
}

type Response struct {
	api.House
}

type CreaterHouse interface {
	CreateHouse(address string, developer string, year int) (api.House, error)
}

func Create(log *slog.Logger, houseCreater CreaterHouse) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userStatus, err := middleware.CheckJWTToken(r)
		if err != nil {
			log.Warn("Invalid token", "err", err)
		}
		if userStatus != api.Moderator {
			log.Info("unauthorized access attempt to house/create", "userStatus", userStatus)
			http.Error(w, "No access", http.StatusUnauthorized)
			return
		}
		var req Request
		err = render.DecodeJSON(r.Body, &req)
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
			http.Error(w, "fail to validate body", http.StatusBadRequest)
			return
		}
		var newHouse api.House
		newHouse, err = houseCreater.CreateHouse(
			req.Address,
			req.Developer,
			req.Year,
		)
		if err != nil {
			log.Error("fail to create house", "err", err)
			err := responses.ServerError(w, r, "fail to create house", 12)
			if err != nil {
				log.Error("fail to send server error code", "err", err)
			}
			return
		}
		log.Info("created house", "newHouse", newHouse)

		render.JSON(w, r, Response{House: newHouse})

	}

}
