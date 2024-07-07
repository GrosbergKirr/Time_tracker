package internal

import (
	"log/slog"
	"os"
)

func SetupLogger() *slog.Logger {
	var log *slog.Logger
	log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	log.Info("Logger set success")
	return log
}
