package house

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"log/slog"
	"main/internal/storage/db"
	"net/http"
)

type SubscriptionRequest struct {
	Email string `json:"email"`
}

func Subscribe(log *slog.Logger, storage *db.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		houseID := vars["id"]

		var req SubscriptionRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Error("Failed to decode request body", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err = storage.SubscribeToHouse(context.Background(), houseID, req.Email)
		if err != nil {
			log.Error("Failed to save subscription", err)
			http.Error(w, "Failed to subscribe", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}
