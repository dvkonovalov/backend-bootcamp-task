package flat

import (
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"log/slog"
	"main/internal/http_server/middleware"
	"main/internal/storage/api"
	"net/http"
)

type RequestCreate struct {
	HouseId int `json:"house_id" validate:"required"`
	Price   int `json:"price" validate:"required"`
	Rooms   int `json:"rooms,omitempty"`
}

type ResponseCreate struct {
	api.Flat
}

type CreatorFlat interface {
	CreateFlat(houseId int, price int, rooms int) (api.Flat, error)
}

func Create(log *slog.Logger, flatCreater CreatorFlat) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userStatus, err := middleware.CheckJWTToken(r)
		if err != nil {
			log.Warn("Invalid token", "err", err)
		}
		if userStatus != api.Client && userStatus != api.Moderator {
			log.Info("unauthorized access attempt to flat/create", "userStatus", userStatus)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		var req RequestCreate
		err = render.DecodeJSON(r.Body, &req)
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
		var newFlat api.Flat
		newFlat, err = flatCreater.CreateFlat(
			req.HouseId,
			req.Price,
			req.Rooms,
		)
		if err != nil {
			log.Error("fail to create flat", "err", err)
			http.Error(w, "fail to create flat", http.StatusInternalServerError)
			return
		}
		log.Info("created flat", "newFlat", newFlat)

		render.JSON(w, r, ResponseCreate{Flat: newFlat})
	}
}
