package domain

import "time"

type Auth struct {
	ID          int64
	UserID      int64
	TgID        int64
	Code        *string
	Token       string
	LastLoginAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}