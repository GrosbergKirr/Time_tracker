package storage

import (
	"log/slog"

	sq "github.com/Masterminds/squirrel"
)

func (s *Storage) DeleteUser(log *slog.Logger, userId int, ok chan bool) error {
	err := s.UserExistenceChecker(log, userId)
	if err != nil {
		log.Error("", slog.Any("err", err))
		return err
	}
	log.Debug("Check user exists success")
	baseQuery := sq.Delete("users").Where(sq.Eq{"id": userId})

	psql := baseQuery.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.ToSql()
	if err != nil {
		log.Error("Failed to build query: ", slog.Any("err", err))
		return err
	}
	stmt, err := s.Db.Prepare(query)
	if err != nil {
		log.Error("Failed to prepare query: ", slog.Any("err", err))
		return err
	}
	log.Debug("Prepare query success")
	_, err = stmt.Query(args...)
	if err != nil {
		log.Error("Failed to execute query: ", slog.Any("err", err))
		return err
	}
	log.Debug("Delete user from db success")
	ok <- true
	return nil
}
