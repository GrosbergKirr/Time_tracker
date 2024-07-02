package api

import (
	"log/slog"
	"net/http"
	"time_track/internal"

	"github.com/go-chi/render"
)

func UserGetter(log *slog.Logger, user UserInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req internal.User

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("fail to decode json", http.StatusBadRequest)
		}
		page := r.URL.Query().Get("page")
		perPage := r.URL.Query().Get("per_page")

		pagination := [2]string{page, perPage}

		res, err := user.GetUser(log, req, pagination)
		if err != nil {
			log.Error("Get data from db mistake", slog.Any("err", err))
			w.WriteHeader(http.StatusInternalServerError)
		}

		render.JSON(w, r, res)
	}
}
