package tools

import (
	"errors"
	"log/slog"
	"regexp"

	"github.com/GrosbergKirr/Time_tracker/internal"
)

func UserValidate(log *slog.Logger, req internal.User, idRequired bool) error {
	forbiddenChars := "!@#$%^&*()-"
	regexpStr := regexp.MustCompile("[" + regexp.QuoteMeta(forbiddenChars) + "]")
	passportLayout := `^\d{4} \d{6}$`
	regexpPassport := regexp.MustCompile(passportLayout)
	err := errors.New("validation error")

	if idRequired {
		if req.Id <= 0 {
			log.Error("Users id is required and should be above 0", slog.Any("err", err))
			return err

		}
	}
	if req.Name != "" {
		if regexpStr.MatchString(req.Name) {
			log.Error("Users data should not constraint !@#$%^&*()-", slog.Any("err", err))
			return err
		}
	}
	if req.Surname != "" {
		if regexpStr.MatchString(req.Surname) {
			log.Error("Users data should not constraint !@#$%^&*()-", slog.Any("err", err))
			return err
		}
	}
	if req.Patronymic != "" {
		if regexpStr.MatchString(req.Patronymic) {
			log.Error("Users data should not constraint !@#$%^&*()-", slog.Any("err", err))
			return err
		}
	}
	if req.Address != "" {
		if regexpStr.MatchString(req.Address) {
			log.Error("Users data should not constraint !@#$%^&*()-", slog.Any("err", err))
			return err
		}
	}
	if req.PassportNum != "" {
		if !regexpPassport.MatchString(req.PassportNum) {
			log.Error("Users data should have layout like 1111 222222", slog.Any("err", err))
			return err
		}
	}
	return nil
}
