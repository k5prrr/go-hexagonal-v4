package service

import (
	"app/internal/app/core/port"
)

type Service struct {
	repo port.IRepo
}

func New(repo port.IRepo) *Service {
	return &Service{repo: repo}
}
