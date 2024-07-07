package api

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/GrosbergKirr/Time_tracker/internal"
	"github.com/GrosbergKirr/Time_tracker/internal/api/client"
	"github.com/GrosbergKirr/Time_tracker/tools"
	"github.com/go-chi/render"
)

// UserCreator godoc
// @Summary Create user
// @Tags users
// @Accept  json
// @Produce  json
// @Param        passportNumber   path      string  true  "Passport data" example({"passportNumber": "1111 123456"})
// @Success 200 "Success"
// @Failure 400 "Invalid input"
// @Failure 500 "Internal server error"
// @Failure 408 "Request timeout"
// @Router /create_user [post]
func UserCreator(ctx context.Context, log *slog.Logger, user UserInterface, clientForSideAPI http.Client, sideApiUrl string) http.HandlerFunc {
	const path string = "api/user_getter"
	return func(w http.ResponseWriter, r *http.Request) {
		var req internal.User
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("fail to decode json", slog.Any("err", err), slog.Any("path", path))
			w.WriteHeader(http.StatusBadRequest)
		}
		log.Info("Get and decode JSON success")

		idIsRequired := false
		if err = tools.UserValidate(log, req, idIsRequired); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Info("Validation true")

		reqToServ := RefactorPasswordForSideAPI(req)
		resp, stat, err := client.GetDataFromSideAPI(log, clientForSideAPI, reqToServ, sideApiUrl)
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
