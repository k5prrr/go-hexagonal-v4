package domain

import (
	"fmt"
	"time"
)

type User struct {
	Uid string

	FamilyName string
	Name       string
	MiddleName string

	BirthDate      time.Time
	Phone          string
	Email          string
	PhoneConfirmed bool
	EmailConfirmed bool

	LastLogin time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time

	Description  string
	PasswordHash string
	KeyApi       string
}

func NewUser() *User {
	return &User{}
}
func (u *User) FullName() string {
	return fmt.Sprintf(
		"%s %s %s",
		u.FamilyName,
		u.Name,
		u.MiddleName,
	)
}
