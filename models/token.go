package models

import "time"

type Token struct {
	UserID    uint
	Token     int
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
