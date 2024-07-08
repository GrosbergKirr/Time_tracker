package api

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/GrosbergKirr/Time_tracker/internal"
	"github.com/go-chi/render"
)

// TaskGetter godoc
// @Summary Get user's tasks
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param   page     query    string     true  "Page number"
// @Param   per_page query    string     true  "Number of items per page"
// @Param   per_page query    string     true  "User id"
// @Success 200 {array} internal.Task "List of tasks"
// @Failure 400 "Invalid input"
// @Failure 500 "Internal server error"
// @Failure 408 "Request timeout"
// @Router /get_users_tasks [get]
func TaskGetter(ctx context.Context, log *slog.Logger, user UserInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const path string = "api/task_getter"
		page := r.URL.Query().Get("page")
		perPage := r.URL.Query().Get("per_page")
		userId := r.URL.Query().Get("user_id")
		log.Info("Get query success")

		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
		ok := make(chan []internal.Task)
		go func() {
			err := user.GetTasks(log, userId, page, perPage, ok)
			if err != nil {
				log.Error("Get data from db error", slog.Any("err: ", err), slog.Any("path", path))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}()
		select {
		case res := <-ok:
			render.JSON(w, r, res)
			w.WriteHeader(http.StatusInternalServerError)
			log.Info("Get tasks success")

			log.Info("Get tasks success")
		case <-ctx.Done():
			log.Error("Timeout error", slog.Any("path", path))
			w.WriteHeader(http.StatusRequestTimeout)

		}

	}
}
