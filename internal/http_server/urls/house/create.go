package house

import (
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"log/slog"
	"main/internal/http_server/middleware"
	"main/internal/storage/api"
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

type HouseCreater interface {
	CreateHouse(address string, developer string, year int) (api.House, error)
}

func Create(log *slog.Logger, houseCreater HouseCreater) http.HandlerFunc {
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
			http.Error(w, "fail to decode body", http.StatusInternalServerError)
			return
		}
		log.Info("request body decoded", slog.Any("req", req))

		err = validator.New().Struct(req)
		if err != nil {
			validatorErr := err.(validator.ValidationErrors)
			log.Error("fail to validate body", "err", validatorErr)
			http.Error(w, "fail to validate body", http.StatusBadRequest)
			return
		}
		var new_house api.House
		new_house, err = houseCreater.CreateHouse(
			req.Address,
			req.Developer,
			req.Year,
		)
		if err != nil {
			log.Error("fail to create house", "err", err)
			http.Error(w, "fail to create house", http.StatusInternalServerError)
			return
		}
		log.Info("created house", "new_house", new_house)

		render.JSON(w, r, Response{House: new_house})

	}

}
