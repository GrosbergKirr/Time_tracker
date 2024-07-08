package storage

import (
	"log/slog"
	"sort"
	"strconv"
	"sync"

	"github.com/GrosbergKirr/Time_tracker/internal"

	sq "github.com/Masterminds/squirrel"
)

func (s *Storage) GetTasks(log *slog.Logger, userId, page, perPage string, ok chan []internal.Task) error {
	const path string = "api/tasks_get"
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		log.Error("convert user id error")
		return err
	}
	err = s.UserExistenceChecker(log, userIdInt)
	if err != nil {
		log.Error("Failed get User", slog.Any("path", path))
		return err
	}
	baseQueryGetTask := sq.Select("*").From("tasks")
	baseQueryGetTask = baseQueryGetTask.Where(sq.Eq{"user_id": userIdInt})

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
	var tasks []internal.Task
	for rows.Next() {
		t := internal.Task{}
		err = rows.Scan(&t.Id, &t.Name, &t.Begin, &t.End, &t.UserId, &t.Status)
		if err != nil {
			log.Error("cant write sql to go-struct", slog.Any("path", path))
			return err
		}
		tasks = append(tasks, t)

		//Сортировка по продолжительности
		sort.Slice(tasks, func(i, j int) bool {
			return (tasks[i].End.Unix() - tasks[i].Begin.Unix()) < (tasks[j].End.Unix() - tasks[j].Begin.Unix())
		})
	}
	log.Debug("Get data from DB success")
	ok <- tasks
	return nil
}
