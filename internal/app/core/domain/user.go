package domain

import (
	"fmt"
	"time"
)

type User struct {
	ID          int64
	FamilyName  string
	Name        string
	MiddleName  string
	Phone       string
	Email       string
	BirthDate   *time.Time `json:"birth_date" gorm:"column:birth_date;type:date"` // может быть NULL
	ParentID    *int64     `json:"parent_id" gorm:"column:parent_id"`             // может быть NULL
	GenderID    *int       `json:"gender_id" gorm:"column:gender_id"`             // может быть NULL
	GroupID     *int64     `json:"group_id" gorm:"column:group_id"`               // может быть NULL
	LastLoginAt *time.Time `json:"last_login_at" gorm:"column:last_login_at"`     // может быть NULL
	Key         string     `json:"key" gorm:"column:key;size:128"`

	CreatedAt time.Time  `json:"created_at" gorm:"column:created_at;not null;autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"column:updated_at;not null;autoUpdateTime"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"column:deleted_at;index"` // soft delete
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
