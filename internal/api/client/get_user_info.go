package client

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/GrosbergKirr/Time_tracker/internal"
	"github.com/go-chi/render"
)

func GetDataFromSideAPI(log *slog.Logger, client http.Client, passport internal.Passport, sideApiUrl string) (internal.User, int, error) {
	body, err := json.Marshal(passport)
	if err != nil {
		log.Error("Marshal json error", err)
	}

	req, err := http.NewRequest("GET", sideApiUrl, bytes.NewBuffer(body))
	if err != nil {
		log.Error("failed to create request", err)
		return internal.User{}, http.StatusInternalServerError, err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Error("failed to get response", err)
		return internal.User{}, resp.StatusCode, err

	}
	defer resp.Body.Close()

	var user internal.User
	err = render.DecodeJSON(resp.Body, &user)
	if err != nil {
		log.Error("fail to decode json", slog.Any("err: ", err))
		return internal.User{}, resp.StatusCode, err
	}
	return user, resp.StatusCode, nil
}
