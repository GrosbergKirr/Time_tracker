package api

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/GrosbergKirr/Time_tracker/internal"
	"github.com/go-chi/render"
)

func UserGetter(ctx context.Context, log *slog.Logger, user UserInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req internal.User
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Error("fail to decode json", slog.Any("err: ", err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		//val, err := tools.ValidatePassport(req.PassportNum)
		//if err != nil || val == false {
		//	log.Error("fail to validate passport number", slog.Any("err: ", err))
		//	w.WriteHeader(http.StatusBadRequest)
		//	return
		//}

		page := r.URL.Query().Get("page")
		perPage := r.URL.Query().Get("per_page")

		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
		ok := make(chan []internal.User)
		go func() {
			err = user.GetUser(log, req, page, perPage, ok)
			if err != nil {
				log.Error("Get data from db error", slog.Any("err", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

		}()
		var res []internal.User
		// Ожидание завершения операции или таймаута
		select {
		case res = <-ok:
			render.JSON(w, r, res)
		case <-ctx.Done():
			w.WriteHeader(http.StatusRequestTimeout)

		}

	}
}
