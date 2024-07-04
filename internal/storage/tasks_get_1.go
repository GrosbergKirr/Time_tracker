package storage

import (
	"log/slog"
	"sort"
	"strconv"

	"github.com/GrosbergKirr/Time_tracker/internal"

	sq "github.com/Masterminds/squirrel"
)

func (s *Storage) GetTasks(log *slog.Logger, user internal.User, page, perPage string, ok chan []internal.Task) error {
	err := s.UserExistenceChecker(log, user.Id)
	if err != nil {
		log.Error("Failed get User", slog.Any("err", err))
		return err
	}

	baseQueryGetTask := sq.Select("*").From("tasks")
	baseQueryGetTask = baseQueryGetTask.Where(sq.Eq{"user_id": user.Id})

	// Пагинация
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		log.Error("Pagination error. Page should be integer")
		return err
	}
	perPageInt, err := strconv.Atoi(perPage)
	if err != nil {
		log.Error("Pagination error. perPage should be integer")
		return err
	}
	offset := (pageInt - 1) * perPageInt
	baseQueryGetTask = baseQueryGetTask.Limit(uint64(perPageInt)).Offset(uint64(offset))

	psql1 := baseQueryGetTask.PlaceholderFormat(sq.Dollar)
	query, args, err := psql1.ToSql()
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

	var tasks []internal.Task
	for rows.Next() {
		t := internal.Task{}
		err = rows.Scan(&t.Id, &t.Name, &t.Begin, &t.End, &t.UserId)
		if err != nil {
			log.Error("cant write sql to go-struct")
			return err
		}
		tasks = append(tasks, t)

		//Сортировка по продолжительности
		sort.Slice(tasks, func(i, j int) bool {
			return (tasks[i].End.Unix() - tasks[i].Begin.Unix()) < (tasks[j].End.Unix() - tasks[j].Begin.Unix())
		})
	}
	ok <- tasks
	return nil
}

func (s *Storage) OneUserGet(log *slog.Logger, passport internal.User) ([]internal.User, error) {
	baseQueryGetUser := sq.Select("*").From("users")

	baseQueryGetUser = baseQueryGetUser.Where(sq.Eq{"passport_number": passport.PassportNum})

	psql := baseQueryGetUser.PlaceholderFormat(sq.Dollar)
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

	var user []internal.User

	for rows.Next() {
		u := internal.User{}
		err = rows.Scan(&u.Id, &u.Name, &u.Surname, &u.Patronymic, &u.Address, &u.PassportNum)
		if err != nil {
			log.Error("cant write sql to go-struct")
			return nil, err
		}
		user = append(user, u)
	}

	return user, nil
}
