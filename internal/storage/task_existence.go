package storage

import (
	"errors"
	"log/slog"

	"github.com/GrosbergKirr/Time_tracker/internal"
	sq "github.com/Masterminds/squirrel"
)

func (s *Storage) TaskExistenceChecker(log *slog.Logger, taskId int) error {
	const path string = "storage/task_existence"
	baseQuery := sq.Select("*").From("tasks").Where(sq.Eq{"id": taskId})

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
	task := internal.Task{}
	for rows.Next() {
		err = rows.Scan(&task.Id, &task.Name, &task.Begin, &task.End, &task.UserId, &task.Status)
		if err != nil {
			log.Error("cant write sql to go-struct", slog.Any("err: ", err), slog.Any("path", path))
			return err
		}
	}
	if task != (internal.Task{}) {
		log.Debug("Task check success")
		return nil
	} else {
		log.Error("User with this id doesn't exists", slog.Any("err: ", err), slog.Any("path", path))
		return errors.New("task doesn't exists")
	}

}
