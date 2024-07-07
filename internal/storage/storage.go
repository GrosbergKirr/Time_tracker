package storage

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

type Storage struct {
	Db *sql.DB
}

func InitStorage(log *slog.Logger, user, pass, addr, name, mode string) *Storage {
	dbPath := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s", user, pass, addr, name, mode)
	db, err := sql.Open("postgres", dbPath)
	if err != nil {
		log.Error("Failed to initialize storage", slog.Any("err", err))
		return nil
	}

	err = goose.Up(db, "migrations")
	if err != nil {
		log.Error("Migration error: %e", slog.Any("err", err))
		return nil
	}
	log.Info("storage initialized")
	return &Storage{Db: db}
}
