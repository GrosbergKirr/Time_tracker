package server

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/GrosbergKirr/Time_tracker/internal"
	"github.com/GrosbergKirr/Time_tracker/internal/api"
	"github.com/GrosbergKirr/Time_tracker/internal/storage"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SetRouters(ctx context.Context, cfg *internal.Config, log *slog.Logger, db *storage.Storage, cli http.Client) chi.Router {
	router := chi.NewRouter()
	router.Post("/get_user_info", api.UserGetter(ctx, log, db))
	router.Get("/get_user_tasks", api.TaskGetter(ctx, log, db))
	router.Post("/create_user", api.UserCreator(ctx, log, db, cli, cfg.ClientUrl))
	router.Post("/make_task", api.TaskMaker(ctx, log, db))
	router.Post("/stop_task", api.TaskStopper(ctx, log, db))
	router.Delete("/delete_user", api.UserDeleter(ctx, log, db))
	router.Post("/update_user", api.UserUpdater(ctx, log, db))

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:9090/swagger/doc.json"), //The url pointing to API definition
	))
	log.Info("Set routers successfully")
	log.Info("USE WEB SWAGGER ON: http://localhost:9090/swagger/index.html#/ ")
	return router
}
