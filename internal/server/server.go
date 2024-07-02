package server

import (
	"log/slog"
	"net/http"
	"time_track/internal"

	"github.com/go-chi/chi/v5"
)

func ServerRun(log *slog.Logger, cfg *internal.Config, router chi.Router) {
	srv := &http.Server{
		Addr:         cfg.Host + cfg.Port,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Error("failed to start server", err)
	} else {
		log.Info("server started")
	}
}
