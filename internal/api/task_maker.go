package api

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/GrosbergKirr/Time_tracker/internal"
	"github.com/GrosbergKirr/Time_tracker/tools"
	"github.com/go-chi/render"
)

// TaskMaker godoc
// @Summary Create task for user
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param        task  body      internal.Task  true  "Task ID" example({"name":"Cook", "user_id":1})
// @Success 200 "Success"
// @Failure 400 "Invalid input"
// @Failure 500 "Internal server error"
// @Failure 408 "Request timeout"
// @Router /make_task [post]
func TaskMaker(ctx context.Context, log *slog.Logger, task UserInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const path string = "api/task_maker"
		var req internal.Task
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("fail to decode json", slog.Any("err: ", err), slog.Any("path", path))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Info("Get and decode JSON success")

		idIsRequired := false
		if err = tools.TaskValidate(log, req, idIsRequired); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Info("Validation true")

		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
		ok := make(chan bool)
		go func() {
			err = task.MakeTask(log, req, ok)
			if err != nil {
				log.Error("Make task error", slog.Any("err: ", err), slog.Any("path", path))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}()
		select {
		case <-ok:
			log.Info("Make task success")
			w.WriteHeader(http.StatusOK)
		case <-ctx.Done():
			log.Error("Timeout error", slog.Any("path", path))
			w.WriteHeader(http.StatusRequestTimeout)

		}
	}
}
