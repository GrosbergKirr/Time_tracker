package api

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/GrosbergKirr/Time_tracker/internal"
	"github.com/GrosbergKirr/Time_tracker/internal/api/client"
	"github.com/GrosbergKirr/Time_tracker/tools"
)

func UserCreator(ctx context.Context, log *slog.Logger, user UserInterface) http.HandlerFunc {
	const path string = "api/user_getter"
	return func(w http.ResponseWriter, r *http.Request) {
		var req internal.User
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Error("fail to decode json", slog.Any("err", err), slog.Any("path", path))
			w.WriteHeader(http.StatusBadRequest)
		}

		val, err := tools.ValidatePassport(req.PassportNum)
		if (err != nil) || val == false {
			log.Error("fail to validate passport number", slog.Any("err: ", err), slog.Any("path", path))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Info("Get and decode JSON success")

		reqToServ := RefactorPasswordForSideAPI(req)
		resp, stat, err := client.GetDataFromSideAPI(log, reqToServ)
		if err != nil {
			log.Error("fail to create client: ", slog.Any("err", err), slog.Any("path", path))
			w.WriteHeader(stat)
			return
		}
		log.Info("Get data from side API success")
		resp.PassportNum = req.PassportNum

		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
		ok := make(chan bool)
		go func() {
			err = user.CreateUser(log, resp, ok)
			if err != nil {
				log.Error("fail to create user: ", slog.Any("err", err), slog.Any("path", path))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}()
		select {
		case <-ok:
			log.Info("User create success")
			w.WriteHeader(http.StatusOK)
		case <-ctx.Done():
			log.Error("Timeout error", slog.Any("path", path))
			w.WriteHeader(http.StatusRequestTimeout)
		}

	}
}

func RefactorPasswordForSideAPI(serNum internal.User) internal.Passport {
	passportSlice := strings.Split(serNum.PassportNum, " ")
	pass := internal.Passport{passportSlice[0], passportSlice[1]}
	return pass
}
