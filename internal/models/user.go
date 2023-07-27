package models

import "time"

type User struct {
	ID        string    `db:"id"`
	Password  string    `db:"password"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Status    bool      `db:"status"`
	CreatedAt time.Time `db:"created_at"`
}
