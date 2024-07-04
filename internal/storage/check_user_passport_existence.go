package storage

import (
	"errors"
	"log/slog"

	"github.com/GrosbergKirr/Time_tracker/internal"

	sq "github.com/Masterminds/squirrel"
)

func (s *Storage) UserPassportExistenceChecker(log *slog.Logger, passport string) error {
	baseQuery := sq.Select("*").From("users").Where(sq.Eq{"passport_number": passport})

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

	rows, err := stmt.Query(args...)
	if err != nil {
		log.Error("Failed to execute query: ", slog.Any("err", err))
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
	if user == (internal.User{}) {
		return nil
	} else {
		log.Error("User with this passport already exists")
		return errors.New("user_exists")
	}

}
