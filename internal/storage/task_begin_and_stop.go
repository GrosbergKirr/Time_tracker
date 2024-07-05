package storage

import (
	"log/slog"
	"time"

	"github.com/GrosbergKirr/Time_tracker/internal"
	sq "github.com/Masterminds/squirrel"
)

func (s *Storage) MakeTask(log *slog.Logger, task internal.Task, ok chan bool) error {
	const path string = "storage/task_begin_and_stop"
	err := s.UserExistenceChecker(log, task.UserId)
	if err != nil {
		log.Error("User existence error ", slog.Any("err: ", err), slog.Any("path", path))
	}
	log.Debug("User exists success")

	baseQuery := sq.Insert("tasks").Columns("name", "time_begin", "user_id").
		Values(task.Name, time.Now(), task.UserId)

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

	_, err = stmt.Query(args...)
	if err != nil {
		log.Error("Failed to execute query: ", slog.Any("path", path))
		return err
	}
	log.Debug("Create task un DB success")
	ok <- true
	return nil
}

func (s *Storage) StopTask(log *slog.Logger, task internal.Task, ok chan bool) error {
	const path string = "storage/task_begin_and_stop"
	err := s.TaskExistenceChecker(log, task.Id)
	if err != nil {
		log.Error("", slog.Any("err: ", err), slog.Any("path", path))
		return err
	}
	log.Debug("Task exists success")
	baseQuery := sq.Update("tasks").Set("time_end", time.Now()).Where(sq.Eq{"id": task.Id})

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

	_, err = stmt.Query(args...)
	if err != nil {
		log.Error("Failed to execute query: ", slog.Any("path", path))
		return err
	}
	log.Debug("Stop Task inn DB success")
	ok <- true
	return nil
}
