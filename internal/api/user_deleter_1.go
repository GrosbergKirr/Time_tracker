package api

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/GrosbergKirr/Time_tracker/internal"
)

func UserDeleter(ctx context.Context, log *slog.Logger, user UserInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const path string = "api/user_deleter"
		var req internal.User
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Error("fail to decode json", slog.Any("err: ", err), slog.Any("path", path))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Info("Get and decode JSON success")

		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
		ok := make(chan bool)
		go func() {
			err = user.DeleteUser(log, req.Id, ok)
			if err != nil {
				log.Error("Delete user error", slog.Any("err: ", err), slog.Any("path", path))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}()

		select {
		case <-ok:
			log.Info("Delete user from success")
			w.WriteHeader(http.StatusOK)
		case <-ctx.Done():
			log.Error("Timeout error", slog.Any("path", path))
			w.WriteHeader(http.StatusRequestTimeout)

		}

	}
}
