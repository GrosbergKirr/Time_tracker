package tools

import (
	"errors"
	"log/slog"
	"regexp"

	"github.com/GrosbergKirr/Time_tracker/internal"
)

func TaskValidate(log *slog.Logger, req internal.Task) error {
	forbiddenChars := "!@#$%^&*()-"
	regexpStr := regexp.MustCompile("[" + regexp.QuoteMeta(forbiddenChars) + "]")
	err := errors.New("validation error")
	if req.Id <= 0 {
		log.Error("task id is required and should be above 0", slog.Any("err", err))
		return err

	}
	if req.Name != "" {
		if regexpStr.MatchString(req.Name) {
			log.Error("Task name should not constraint !@#$%^&*()-", slog.Any("err", err))
			return err
		}
	}
	if req.UserId <= 0 {
		if regexpStr.MatchString(req.Name) {
			log.Error("Users id is required and should be above 0", slog.Any("err", err))
			return err
		}
	}
	return nil
}
