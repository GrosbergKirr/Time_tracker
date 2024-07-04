package storage

import (
	"log/slog"
	"sync"

	"github.com/GrosbergKirr/Time_tracker/internal"
	sq "github.com/Masterminds/squirrel"
)

func (s *Storage) UpdateUser(log *slog.Logger, user internal.User, ok chan bool) error {
	const path string = "storage/user_updater"
	err := s.UserExistenceChecker(log, user.Id)
	if err != nil {
		log.Error("user existence error", slog.Any("err: ", err), slog.Any("path", path))
		return err
	}
	log.Debug("Check user id is valid success", slog.Any("id", user.Id))
	baseQuery := sq.Update("users")
	if user.Name != "" {
		baseQuery = baseQuery.Set("name", user.Name)
	}
	if user.Surname != "" {
		baseQuery = baseQuery.Set("surname", user.Surname)
	}
	if user.Patronymic != "" {
		baseQuery = baseQuery.Set("patronymic", user.Patronymic)
	}
	if user.Address != "" {
		baseQuery = baseQuery.Set("address", user.Address)
	}
	if user.PassportNum != "" {
		baseQuery = baseQuery.Set("passport_number", user.PassportNum)
	}
	baseQuery = baseQuery.Where(sq.Eq{"id": user.Id})

	psql := baseQuery.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.ToSql()
	if err != nil {
		log.Error("Failed to build query: ", slog.Any("err: ", err), slog.Any("path", path))
		return err
	}
	stmt, err := s.Db.Prepare(query)
	if err != nil {
		log.Error("Failed to prepare query: ", slog.Any("err: ", err), slog.Any("path", path))
		return err
	}
	log.Debug("Prepare query success")
	mu := sync.Mutex{}
	mu.Lock()
	_, err = stmt.Query(args...)
	if err != nil {
		log.Error("Failed to execute query: ", slog.Any("err: ", err), slog.Any("path", path))
		return err
	}
	mu.Unlock()
	log.Debug("Update success")
	ok <- true
	return nil

}
