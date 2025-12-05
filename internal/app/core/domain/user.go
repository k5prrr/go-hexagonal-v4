package domain

import (
	"fmt"
	"time"
)

type User struct {
	ID         int64
	FamilyName *string
	Name       *string
	MiddleName *string

	Phone     string
	Email     *string
	BirthDate *time.Time `json:"birth_date" gorm:"column:birth_date;type:date"` // может быть NULL

	ParentID *int64 `json:"parent_id" gorm:"column:parent_id"` // может быть NULL
	GenderID *int   `json:"gender_id" gorm:"column:gender_id"` // может быть NULL
	RoleID   *int64 `json:"role_id" gorm:"column:role_id"`     // может быть NULL

	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;not null;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;not null;autoUpdateTime"`
}

func (u *User) FullName() string {
	return fmt.Sprintf(
		"%s %s %s",
		u.FamilyName,
		u.Name,
		u.MiddleName,
	)
}

type UserGroup struct {
	ID   int64
	Name string
}

/*
(1,	'Админ'),
(2,	'Менеджер'),
(3,	'Продавец'),
(4,	'Новый');
*/

type Auth struct {
	ID          int64
	UserID      int64
	TgID        int64
	Code        *string
	Secret      string
	LastLoginAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UserFull struct {
	User *User
	Auth *Auth
}
