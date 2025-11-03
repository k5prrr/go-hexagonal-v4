package domain

import "time"

type CheckPhone struct {
	UUID string
}
type AuthCode struct {
	ID        int64
	Code      string
	Type      string
	UUID      *string
	Phone     *string
	CreatedAt time.Time
	UpdatedAt time.Time
}
