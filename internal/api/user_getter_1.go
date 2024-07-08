package api

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/GrosbergKirr/Time_tracker/internal"
	"github.com/GrosbergKirr/Time_tracker/tools"
	"github.com/go-chi/render"
)

// UserGetter godoc
// @Summary Get users
// @Tags users
// @Accept  json
// @Produce  json
// @Param   page     query    string     true  "Page number"
// @Param   per_page query    string     true  "Number of items per page"
// @Param   user     body     internal.User false "User details" example(internal.User){"id": 5, "name":"Ivan", "surname":"Ivanov", "patronymic": "Ivanovich", "address": "SPB", "passportNumber":"1111 123456"}
// @Success 200 {array} internal.User "List of users"
// @Failure 400 "Invalid input"
// @Failure 500 "Internal server error"
// @Failure 408 "Request timeout"
// @Router /get_user_info [post]
func UserGetter(ctx context.Context, log *slog.Logger, user UserInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const path string = "api/user_getter"
		var req internal.User
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("fail to decode json", slog.Any("err: ", err), slog.Any("path", path))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		page := r.URL.Query().Get("page")
		perPage := r.URL.Query().Get("per_page")
		log.Info("Get and decode JSON success")

		idIsRequired := false
		if err = tools.UserValidate(log, req, idIsRequired); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Info("Validation true")

		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
		ok := make(chan []internal.User)
		go func() {
			err = user.GetUser(log, req, page, perPage, ok)
			if err != nil {
				log.Error("Get data from db error", slog.Any("err: ", err), slog.Any("path", path))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

		}()
		var res []internal.User
		select {
		case res = <-ok:
			if render.JSON(w, r, res); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			log.Info("Get user success")
		case <-ctx.Done():
			log.Error("Timeout error", slog.Any("path", path))
			w.WriteHeader(http.StatusRequestTimeout)

		}

	}
}
