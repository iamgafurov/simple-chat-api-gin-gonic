package models

import (
	"time"
)

type User struct {
	ID      int64     `json:"id"`
	Name    string    `json:"name"`
	Login   string    `json:"login"`
	Created time.Time `json:"created"`
}

type Registration struct {
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type ID struct {
	ID int64 `uri:"id"`
}

type Auth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
