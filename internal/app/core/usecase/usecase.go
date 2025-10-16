package usecase

import "app/internal/app/core/port"

type Usecase struct {
	repo *port.IAuthRepo
}

func NewUuth(repo port.IAuthRepo) *Auth {
	return &Auth{repo: repo}
}
