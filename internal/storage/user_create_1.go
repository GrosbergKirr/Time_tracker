package storage

import (
	"log/slog"
	"sync"

	"github.com/GrosbergKirr/Time_tracker/internal"

	sq "github.com/Masterminds/squirrel"
)

func (s *Storage) CreateUser(log *slog.Logger, user internal.User, ok chan bool) error {
	const path string = "storage/user_create"
	err := s.UserPassportExistenceChecker(log, user.PassportNum)
	if err != nil {
		log.Error("Failed to build query: ", slog.Any("err", err), slog.Any("path", path))
	}
	log.Debug("Check passport uniq success")
	baseQuery := sq.Insert("users").Columns("name", "surname", "patronymic", "address", "passport_number").
		Values(user.Name, user.Surname, user.Patronymic, user.Address, user.PassportNum)

	psql := baseQuery.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.ToSql()
	if err != nil {
		log.Error("Failed to build query: ", slog.Any("err", err), slog.Any("path", path))
		return err
	}
	stmt, err := s.Db.Prepare(query)
	if err != nil {
		log.Error("Failed to prepare query: ", slog.Any("err", err), slog.Any("path", path))
		return err
	}
	log.Debug("Prepare query success")
	mu := sync.Mutex{}
	mu.Lock()
	_, err = stmt.Query(args...)
	if err != nil {
		log.Error("Failed to execute query: ", slog.Any("err", err), slog.Any("path", path))
		return err
	}
	mu.Unlock()
	log.Debug("create user in DB success")
	ok <- true
	return nil
}
