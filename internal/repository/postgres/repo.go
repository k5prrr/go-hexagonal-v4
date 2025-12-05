package postgres

import (
	"app/pkg/database"
	"errors"
)

type Repo struct {
	db database.IDB
}

func New(db database.IDB) *Repo {
	return &Repo{db: db}
}

var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrUnsupported  = errors.New("unsupported filter key")
)
