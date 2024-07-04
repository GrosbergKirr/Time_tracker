package storage

import (
	"errors"
	"log/slog"

	"github.com/GrosbergKirr/Time_tracker/internal"
	sq "github.com/Masterminds/squirrel"
)

func (s *Storage) UserExistenceChecker(log *slog.Logger, userId int) error {
	const path string = "storage/user_existence_checker"
	baseQuery := sq.Select("*").From("users").Where(sq.Eq{"id": userId})

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
	rows, err := stmt.Query(args...)
	if err != nil {
		log.Error("Failed to execute query: ", slog.Any("err: ", err), slog.Any("path", path))
		return err
	}
	user := internal.User{}
	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Name, &user.Surname, &user.Patronymic, &user.Address, &user.PassportNum)
		if err != nil {
			log.Error("cant write sql to go-struct")
			return err
		}
	}
	if user != (internal.User{}) {
		log.Debug("User check success")
		return nil
	} else {
		log.Error("User with this id doesn't exists", slog.Any("err: ", err), slog.Any("path", path))
		return errors.New("user doesn't exists")
	}

}
