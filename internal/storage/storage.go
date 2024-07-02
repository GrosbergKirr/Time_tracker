package storage

import (
	"database/sql"
	"fmt"
	"os"

	//_ "github.com/jackc/pgx/v5/stdlib"
	//"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

type Storage struct {
	Db *sql.DB
}

//	func GetStandardDB(sqlxDB *sqlx.DB) *sql.DB {
//		return sqlxDB.DB
//	}
func InitStorage(user, pass, name, mode string) (*Storage, error) {
	dbPath := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s?sslmode=%s", user, pass, name, mode)
	db, err := sql.Open("postgres", dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	//stdDb := GetStandardDB(db)
	err = goose.Up(db, "migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to up migration %e", err)
	}
	return &Storage{Db: db}, nil
}
