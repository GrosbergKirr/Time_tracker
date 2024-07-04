package api

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/GrosbergKirr/Time_tracker/internal"
)

func TaskGetter(ctx context.Context, log *slog.Logger, user UserInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const path string = "api/task_getter"
		var req internal.User
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Error("fail to decode json", slog.Any("err: ", err), slog.Any("path", path))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		page := r.URL.Query().Get("page")
		perPage := r.URL.Query().Get("per_page")
		log.Info("Get and decode JSON success")

		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
		ok := make(chan []internal.Task)
		go func() {
			err := user.GetTasks(log, req, page, perPage, ok)
			if err != nil {
				log.Error("Get data from db error", slog.Any("err: ", err), slog.Any("path", path))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}()
		select {
		case res := <-ok:

			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			log.Info("Get tasks success")
		case <-ctx.Done():
			log.Error("Timeout error", slog.Any("path", path))
			w.WriteHeader(http.StatusRequestTimeout)

		}

	}
}
