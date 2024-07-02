package api

import (
	"log/slog"
	"time_track/internal"
)

type UserInterface interface {
	GetUser(log *slog.Logger, user internal.User, pagination [2]string) ([]internal.User, error)
}
