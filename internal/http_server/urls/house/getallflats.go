package house

import (
	"github.com/go-chi/render"
	"github.com/gorilla/mux"
	"log/slog"
	"main/internal/storage/api"
	"net/http"
	"strconv"
)

type ResponseGetAllFlats struct {
	Flats []api.Flat `json:"flats"`
}

type AllFlatsGetter interface {
	GetAllFlats(house_id int) ([]api.Flat, error)
}

func GetFlats(log *slog.Logger, allFlatsGetter AllFlatsGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Error("Invalid ID", "id", idStr, "error", err)
			return
		}

		var flats []api.Flat
		flats, err = allFlatsGetter.GetAllFlats(id)
		if err != nil {
			log.Error("fail to get all flats", "err", err)
			return
		}
		log.Info("Get all flat in house", "flats", flats)
		render.JSON(w, r, ResponseGetAllFlats{Flats: flats})
	}
}
