package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/GrosbergKirr/Time_tracker/docs"
	"github.com/GrosbergKirr/Time_tracker/internal"
	"github.com/GrosbergKirr/Time_tracker/internal/api"
	"github.com/GrosbergKirr/Time_tracker/internal/server"
	"github.com/GrosbergKirr/Time_tracker/internal/storage"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:9090

func main() {
	log := internal.SetupLogger()
	cfg := internal.SetupConfig(log)
	db := storage.InitStorage(log, cfg.Username, cfg.Password, cfg.Adress, cfg.Database, cfg.Mode)

	ctx := context.Background()

	client := http.Client{}
	router := chi.NewRouter()

	router.Get("/get_user_info", api.UserGetter(ctx, log, db))
	router.Get("/get_user_tasks", api.TaskGetter(ctx, log, db))
	router.Post("/create_user", api.UserCreator(ctx, log, db, client, cfg.ClientUrl))
	router.Post("/make_task", api.TaskMaker(ctx, log, db))
	router.Post("/stop_task", api.TaskStopper(ctx, log, db))
	router.Delete("/delete_user", api.UserDeleter(ctx, log, db))
	router.Post("/update_user", api.UserUpdater(ctx, log, db))

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:9090/swagger/doc.json"), //The url pointing to API definition
	))

	log.Info("USE WEB SWAGGER ON: http://localhost:9090/swagger/index.html#/ ")

	serverStopSig := make(chan os.Signal)
	newServer := server.NewServer(cfg, router)
	go newServer.ServerRun(log, cfg)

	signal.Notify(serverStopSig, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	<-serverStopSig
	newServer.ServerStop(ctx, log)

}
