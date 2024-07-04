package main

import (
	"context"
	"log/slog"

	"github.com/GrosbergKirr/Time_tracker/internal"
	"github.com/GrosbergKirr/Time_tracker/internal/api"
	"github.com/GrosbergKirr/Time_tracker/internal/server"
	"github.com/GrosbergKirr/Time_tracker/internal/storage"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/go-chi/chi/v5"

	_ "github.com/GrosbergKirr/Time_tracker/docs"
)

func main() {
	log := internal.SetupLogger()
	cfg := internal.SetupConfig(log)

	log.Info("config % logger set success")
	db, err := storage.InitStorage(cfg.Username, cfg.Password, cfg.Database, cfg.Mode)
	if err != nil {
		log.Info("Failed to initialize storage")
	}
	log.Info("storage initialized")

	ctx := context.Background()
	router := chi.NewRouter()

	router.Get("/get_user_info", api.UserGetter(ctx, log, db))
	router.Get("/get_users_tasks", api.TaskGetter(ctx, log, db))
	router.Post("/create_user", api.UserCreator(ctx, log, db))
	router.Post("/make_task", api.TaskMaker(ctx, log, db))
	router.Post("/stop_task", api.TaskStopper(ctx, log, db))
	router.Post("/delete_user", api.UserDeleter(ctx, log, db))
	router.Post("/update_user", api.UserUpdater(ctx, log, db))

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:9090/swagger/doc.json"), //The url pointing to API definition
	))

	log.Info("starting server", slog.String("address", cfg.Host+cfg.Port))

	//Run server
	server.ServerRun(log, cfg, router)

}
