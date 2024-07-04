package client

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/GrosbergKirr/Time_tracker/internal"
)

func GetDataFromSideAPI(log *slog.Logger, client http.Client, passport internal.Passport) (internal.User, int, error) {

	urlBody, exists := os.LookupEnv("CLIENT_URL")
	if !exists {
		log.Error("set CLIENT_URL env variable")
		return internal.User{}, http.StatusInternalServerError, nil
	}

	body, err := json.Marshal(passport)
	if err != nil {
		log.Error("Marshal json error", err)
	}
	req, err := http.NewRequest("GET", urlBody, bytes.NewBuffer(body))
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

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("failed to read response body", err)
		return internal.User{}, resp.StatusCode, err
	}
	var user internal.User
	err = json.Unmarshal(responseBody, &user)
	if err != nil {
		log.Error("failed to unmarshal response body", err)
		return internal.User{}, resp.StatusCode, err
	}
	return user, resp.StatusCode, nil
}
