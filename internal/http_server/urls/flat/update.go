package flat

import (
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"log/slog"
	"main/internal/storage/api"
	"net/http"
)

type Request_update struct {
	Id     int    `json:"id" validate:"required"`
	Status string `json:"status"`
}

type Response_update struct {
	api.Flat
}

type FlatUpdater interface {
	UpdateFlat(id int, status string) (api.Flat, error)
}

func Update(log *slog.Logger, flatUpdater FlatUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request_update
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
		var update_flat api.Flat
		update_flat, err = flatUpdater.UpdateFlat(
			req.Id,
			req.Status,
		)
		if err != nil {
			log.Error("fail to update flat", "err", err)
			render.JSON(w, r, "failed to update flat")
			return
		}
		log.Info("update flat", "update_flat", update_flat)

		render.JSON(w, r, Response_update{Flat: update_flat})
	}
}
