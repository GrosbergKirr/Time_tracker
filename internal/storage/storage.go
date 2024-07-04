package storage

import (
	"database/sql"
	"fmt"
	"log"

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
		log.Fatalf("Storage init error: %e", err)
	}
	//stdDb := GetStandardDB(db)
	err = goose.Up(db, "migrations")
	if err != nil {
		log.Fatalf("Migration error error: %e", err)
	}
	return &Storage{Db: db}, nil
}
