package storage

import (
	"log/slog"
	"strconv"
	"sync"

	"github.com/GrosbergKirr/Time_tracker/internal"

	sq "github.com/Masterminds/squirrel"
)

func (s *Storage) GetUser(log *slog.Logger, user internal.User, page string, perPage string, ok chan []internal.User) error {
	const path string = "path = storage/user_get"
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
	if user.PassportNum != "" {
		baseQuery = baseQuery.Where(sq.Eq{"pasport_number": user.PassportNum})
	}

	// Пагинация
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		log.Error("Pagination error. Page should be integer", slog.Any("path", path))
		return err
	}
	perPageInt, err := strconv.Atoi(perPage)
	if err != nil {
		log.Error("Pagination error. perPage should be integer", slog.Any("path", path))
		return err
	}
	offset := (pageInt - 1) * perPageInt
	baseQuery = baseQuery.Limit(uint64(perPageInt)).Offset(uint64(offset))

	psql := baseQuery.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.ToSql()
	if err != nil {
		log.Error("Failed to build query: ", slog.Any("path", path))
		return err
	}
	stmt, err := s.Db.Prepare(query)
	if err != nil {
		log.Error("Failed to prepare query: ", slog.Any("path", path))
		return err
	}
	log.Debug("Prepare query success")
	mu := sync.RWMutex{}
	mu.Lock()
	rows, err := stmt.Query(args...)
	if err != nil {
		log.Error("Failed to execute query: ", slog.Any("path", path))
		return err
	}
	mu.Unlock()
	var res []internal.User
	for rows.Next() {
		r := internal.User{}
		err = rows.Scan(&r.Id, &r.Name, &r.Surname, &r.Patronymic, &r.Address, &r.PassportNum)
		if err != nil {
			log.Error("cant write sql to go-struct", slog.Any("path", path))
			return err
		}
		res = append(res, r)

	}
	log.Debug("Get data from DB success")
	ok <- res
	return nil
}
