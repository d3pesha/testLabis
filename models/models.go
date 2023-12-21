package models

import "time"

type Client struct {
	ID       int      `json:"id"`
	FIO      string   `json:"fio"`
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Objects  []Object `json:"objects"`
}

type Object struct {
	ID        int        `json:"id" db:"id"`
	UserID    int        `json:"id_user" db:"id_user"`
	Address   string     `json:"address" db:"address"`
	IsVisible bool       `json:"is_visible" db:"is_visible"`
	Contracts []Contract `json:"contracts"`
}

type Contract struct {
	ID       int       `json:"id" db:"id"`
	ObjectID int       `json:"object_id" db:"id_object"`
	Data     time.Time `json:"data" db:"data"`
	Number   string    `json:"number" db:"number"`
	Status   bool      `json:"status" db:"status"`
}

type ErrorResponse struct {
	Message string `json:"error"`
}
