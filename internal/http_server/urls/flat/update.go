package flat

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

type RequestUpdate struct {
	Id     int    `json:"id" validate:"required"`
	Status string `json:"status"`
}

type ResponseUpdate struct {
	api.Flat
}

type UpdaterFlat interface {
	GetStatus(id int) (string, error)
	UpdateFlat(id int, status string, moderator string) (api.Flat, error)
	BlockModerationOtherAdmin(flatId int, moderator string) (bool, error)
}

func Update(log *slog.Logger, flatUpdater UpdaterFlat) http.HandlerFunc {
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
		var req RequestUpdate
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

		usernameAdmin, err := middleware.CheckGetUser(r)
		if err != nil {
			log.Error("fail to check user in flat/update", "err", err)
			err := responses.ServerError(w, r, "fail to check user in flat/update", 21)
			if err != nil {
				log.Error("fail to send server error code", "err", err)
			}
			return
		}

		var updateFlat api.Flat
		updateFlat, err = flatUpdater.UpdateFlat(
			req.Id,
			req.Status,
			usernameAdmin,
		)
		if err != nil {
			log.Error("fail to update flat", "err", err)
			http.Error(w, "This apartment is already being moderated", http.StatusForbidden)
			return
		}
		log.Info("update flat", "updateFlat", updateFlat)

		res, err := flatUpdater.BlockModerationOtherAdmin(req.Id, usernameAdmin)
		if err != nil {
			log.Error("fail to flat/update block flat with moderator", "err", err)
			err := responses.ServerError(w, r, "fail to flat/update block flat with moderator", 22)
			if err != nil {
				log.Error("fail to send server error code", "err", err)
			}
			return
		}
		if res == false {
			log.Error("fail to flat/update block flat with moderator", "res", res)
			err := responses.ServerError(w, r, "fail to flat/update block flat with moderator", 22)
			if err != nil {
				log.Error("fail to send server error code", "err", err)
			}
			return
		}

		render.JSON(w, r, ResponseUpdate{Flat: updateFlat})
	}
}
