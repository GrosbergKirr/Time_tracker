package storage

import (
	"log/slog"
	"strconv"
	"time_track/internal"

	sq "github.com/Masterminds/squirrel"
	//_ "github.com/jackc/pgx/v5/stdlib"
	//_ "github.com/jmoiron/sqlx"
)

func (s *Storage) GetUser(log *slog.Logger, user internal.User, pagination [2]string) ([]internal.User, error) {

	baseQuery := sq.Select("*").From("users")

	if user.Id != 0 {
		baseQuery = baseQuery.Where(sq.Eq{"id": user.Id})
	}
	if user.Name != "" {
		baseQuery = baseQuery.Where(sq.Eq{"name": user.Name})
	}
	if user.Surname != "" {
		baseQuery = baseQuery.Where(sq.Eq{"surname": user.Surname})
	}
	if user.Patronymic != "" {
		baseQuery = baseQuery.Where(sq.Eq{"patronymic": user.Patronymic})
	}
	if user.Address != "" {
		baseQuery = baseQuery.Where(sq.Eq{"address": user.Address})
	}
	if user.Seria != 0 {
		baseQuery = baseQuery.Where(sq.Eq{"pasport_seria": user.Seria})
	}
	if user.Num != 0 {
		baseQuery = baseQuery.Where(sq.Eq{"pasport_num": user.Num})
	}

	// Пагинация
	page, _ := strconv.Atoi(pagination[0])
	perPage, err := strconv.Atoi(pagination[1])
	offset := (page - 1) * perPage
	baseQuery = baseQuery.Limit(uint64(perPage)).Offset(uint64(offset))

	psql := baseQuery.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.ToSql()

	if err != nil {
		log.Error("Failed to build query: ", slog.Any("err", err))
		return nil, err
	}
	stmt, err := s.Db.Prepare(query)
	if err != nil {
		log.Error("Failed to prepare query: ", slog.Any("err", err))
		return nil, err
	}

	rows, err := stmt.Query(args...)
	if err != nil {
		log.Error("Failed to execute query: ", slog.Any("err", err))
		return nil, err
	}
	//
	var res []internal.User
	for rows.Next() {
		r := internal.User{}
		err = rows.Scan(&r.Id, &r.Name, &r.Surname, &r.Patronymic, &r.Address, &r.Seria, &r.Num)
		if err != nil {
			log.Error("cant write sql to go-struct")
			return nil, err
		}
		res = append(res, r)

	}
	return res, nil
}
