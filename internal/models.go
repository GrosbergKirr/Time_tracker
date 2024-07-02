package internal

import (
	"time"
)

type User struct {
	Id         int    `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	Surname    string `json:"surname,omitempty"`
	Patronymic string `json:"patronymic,omitempty"`
	Address    string `json:"address,omitempty"`
	Seria      int    `json:"pasport_seria,omitempty"`
	Num        int    `json:"pasport_num,omitempty"`
}

type Task struct {
	Id    int
	Begin time.Time
	End   time.Time
	User  User
}
