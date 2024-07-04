package api

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/GrosbergKirr/Time_tracker/internal"
)

func TaskStopper(ctx context.Context, log *slog.Logger, task UserInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const path string = "api/task_stopper"
		var req internal.Task
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
			err = task.StopTask(log, req, ok)
			if err != nil {
				log.Error("Stop task error", slog.Any("err: ", err), slog.Any("path", path))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}()

		select {
		case <-ok:
			log.Info("Stop task success")
			w.WriteHeader(http.StatusOK)
		case <-ctx.Done():
			log.Error("Timeout error", slog.Any("path", path))
			w.WriteHeader(http.StatusRequestTimeout)

		}
	}
}
