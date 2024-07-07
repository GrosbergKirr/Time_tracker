package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/GrosbergKirr/Time_tracker/docs"
	"github.com/GrosbergKirr/Time_tracker/internal"
	"github.com/GrosbergKirr/Time_tracker/internal/server"
	"github.com/GrosbergKirr/Time_tracker/internal/storage"
)

// @title Time Tracker app swagger
// @version 1.0
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
	db := storage.InitStorage(log, cfg.Username, cfg.Password, cfg.Address, cfg.Database, cfg.Mode)

	ctx := context.Background()

	client := http.Client{}

	router := server.SetRouters(ctx, cfg, log, db, client)

	serverStopSig := make(chan os.Signal)
	newServer := server.NewServer(cfg, router)
	go newServer.ServerRun(log, cfg)

	signal.Notify(serverStopSig, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	<-serverStopSig
	newServer.ServerStop(ctx, log)
}
