package main

import (
	"log/slog"
	"time_track/internal"
	"time_track/internal/api"
	"time_track/internal/server"
	"time_track/internal/storage"

	"github.com/go-chi/chi/v5"
)

func main() {

	cfg := internal.SetupConfig()
	log := internal.SetupLogger()

	log.Info("config % logger set success")
	storage, err := storage.InitStorage(cfg.Username, cfg.Password, cfg.Database, cfg.Mode)
	if err != nil {
		log.Info("Failed to initialize storage")
	}
	_ = storage

	log.Info("storage initialized")
	router := chi.NewRouter()

	//Set handel's paths
	router.Get("/get_user", api.UserGetter(log, storage))

	log.Info("starting server", slog.String("address", cfg.Host+cfg.Port))

	//Run server
	server.ServerRun(log, cfg, router)

}
