package api

import (
	"log/slog"

	"github.com/GrosbergKirr/Time_tracker/internal"
)

type UserInterface interface {
	GetUser(log *slog.Logger, user internal.User, page string, perPage string, ok chan []internal.User) error
	GetTasks(log *slog.Logger, userId, page, perPage string, ok chan []internal.Task) error
	CreateUser(log *slog.Logger, user internal.User, ok chan bool) error
	MakeTask(log *slog.Logger, task internal.Task, ok chan bool) error
	StopTask(log *slog.Logger, task internal.Task, ok chan bool) error
	DeleteUser(log *slog.Logger, userId int, ok chan bool) error
	UpdateUser(log *slog.Logger, user internal.User, ok chan bool) error
}
