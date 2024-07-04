package storage

import (
	"database/sql"
	"fmt"
	"log/slog"

	//_ "github.com/jackc/pgx/v5/stdlib"
	//"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

type Storage struct {
	Db *sql.DB
}

func InitStorage(log *slog.Logger, user, pass, addr, name, mode string) (*Storage, error) {
	dbPath := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s", user, pass, addr, name, mode)
	db, err := sql.Open("postgres", dbPath)
	if err != nil {
		log.Error("Storage init error: %e", err)
		return nil, err
	}

	err = goose.Up(db, "migrations")
	if err != nil {
		log.Error("Migration error: %e", err)
		return nil, err
	}
	return &Storage{Db: db}, nil
}
