package dto

import "time"

type RegisterReqDto struct {
	ID       string `json:"id"`
	Email    string `json:"email" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginReqDto struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResDto struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Status    bool      `json:"status"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
}
