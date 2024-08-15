package flat

import (
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"log/slog"
	"main/internal/storage/api"
	"net/http"
)

type Request_create struct {
	House_id int `json:"house_id" validate:"required"`
	Price    int `json:"price" validate:"required"`
	Rooms    int `json:"rooms,omitempty"`
}

type Response_create struct {
	api.Flat
}

type FlatCreator interface {
	CreateFlat(house_id int, price int, rooms int) (api.Flat, error)
}

func Create(log *slog.Logger, flatCreater FlatCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request_create
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
		var new_flat api.Flat
		new_flat, err = flatCreater.CreateFlat(
			req.House_id,
			req.Price,
			req.Rooms,
		)
		if err != nil {
			log.Error("fail to create flat", "err", err)
			render.JSON(w, r, "failed to create flat")
			return
		}
		log.Info("created flat", "new_flat", new_flat)

		render.JSON(w, r, Response_create{Flat: new_flat})
	}
}
