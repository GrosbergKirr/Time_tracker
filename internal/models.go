package internal

import (
	"time"
)

type User struct {
	Id          int    `json:"id,omitempty"`
	Surname     string `json:"surname,omitempty"`
	Name        string `json:"name,omitempty"`
	Patronymic  string `json:"patronymic,omitempty"`
	Address     string `json:"address,omitempty"`
	PassportNum string `json:"passportNumber,omitempty"`
}

type Task struct {
	Id     int       `json:"id,omitempty"`
	Name   string    `json:"name,omitempty"`
	Begin  time.Time `json:"time_begin,omitempty"`
	End    time.Time `json:"time_end,omitempty"`
	UserId int       `json:"user_id,omitempty"`
}

type Passport struct {
	PassportSerie  string `json:"passportSerie"`
	PassportNumber string `json:"passportNumber"`
}
