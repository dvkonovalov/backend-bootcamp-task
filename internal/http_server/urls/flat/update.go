package flat

import (
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"log/slog"
	"main/internal/http_server/middleware"
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
	GetStatus(id int) (string, error)
	UpdateFlat(id int, status string) (api.Flat, error)
}

func Update(log *slog.Logger, flatUpdater FlatUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userStatus, err := middleware.CheckJWTToken(r)
		if err != nil {
			log.Warn("Invalid token", "err", err)
		}
		if userStatus != api.Moderator {
			log.Info("unauthorized access attempt to house/update", "userStatus", userStatus)
			http.Error(w, "No access", http.StatusUnauthorized)
			return
		}
		var req Request_update
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
		var update_flat api.Flat
		update_flat, err = flatUpdater.UpdateFlat(
			req.Id,
			req.Status,
		)
		if err != nil {
			log.Error("fail to update flat", "err", err)
			http.Error(w, "fail to update flat", http.StatusInternalServerError)
			return
		}
		log.Info("update flat", "update_flat", update_flat)

		render.JSON(w, r, Response_update{Flat: update_flat})
	}
}
