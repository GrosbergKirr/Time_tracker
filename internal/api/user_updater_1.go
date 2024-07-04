package api

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/GrosbergKirr/Time_tracker/internal"
	"github.com/GrosbergKirr/Time_tracker/tools"
)

func UserUpdater(ctx context.Context, log *slog.Logger, user UserInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const path string = "api/user_updater"
		var req internal.User
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Error("fail to decode json", slog.Any("err: ", err), slog.Any("path", path))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if req.PassportNum != "" {
			val, err := tools.ValidatePassport(req.PassportNum)
			if (err != nil) || val == false {
				log.Error("fail to validate passport number", slog.Any("err: ", err), slog.Any("path", path))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		log.Info("Get and decode JSON success")
		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
		ok := make(chan bool)
		go func() {
			err = user.UpdateUser(log, req, ok)
			if err != nil {
				log.Error("fail to update user data", slog.Any("err: ", err), slog.Any("path", path))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}()
		select {
		case <-ok:
			log.Info("Get tasks from db success")
			w.WriteHeader(http.StatusOK)
		case <-ctx.Done():
			w.WriteHeader(http.StatusRequestTimeout)
			log.Error("Timeout error", slog.Any("path", path))

		}

	}
}
