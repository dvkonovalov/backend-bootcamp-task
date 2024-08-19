package house

import (
	"github.com/go-chi/render"
	"github.com/gorilla/mux"
	"log/slog"
	"main/internal/http_server/middleware"
	"main/internal/storage/api"
	"net/http"
	"strconv"
)

type ResponseGetAllFlats struct {
	Flats []api.Flat `json:"flats"`
}

type AllFlatsGetter interface {
	GetAllFlats(houseId int, userType string) ([]api.Flat, error)
}

func GetFlats(log *slog.Logger, allFlatsGetter AllFlatsGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userType, err := middleware.CheckJWTToken(r)
		if err != nil {
			log.Warn("Invalid token", "err", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		if userType != api.Client && userType != api.Moderator {
			log.Info("unauthorized access attempt to house/{id}")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		vars := mux.Vars(r)
		idStr := vars["id"]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Error("Invalid ID", "id", idStr, "error", err)
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		var flats []api.Flat
		flats, err = allFlatsGetter.GetAllFlats(id, userType)
		if err != nil {
			log.Error("fail to get all flats", "err", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		log.Info("Get all flat in house", "flats", flats)
		render.JSON(w, r, ResponseGetAllFlats{Flats: flats})
	}
}
