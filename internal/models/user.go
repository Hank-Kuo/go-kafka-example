package models

import "time"

type User struct {
	ID        string    `json:"id" db:"id"`
	Account   string    `json:"account" db:"account"`
	Password  string    `json:"password" db:"password"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Status    bool      `json:"status" db:"status"`
	CreatedAt time.Time `json:"create_at" db:"created_at"`
}
