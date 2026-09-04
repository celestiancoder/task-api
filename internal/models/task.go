package models

import "time"

type Task struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TaskQueryParams struct {
	Completed string
	Search    string
	Limit     int
	Page      int
	Sort      string
	Order     string
}