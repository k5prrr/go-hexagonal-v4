package usecase

import (
	"app/internal/app/core/port"
	"app/internal/app/core/service"
)

// Тут он как основной, тут только создания
type UseCase struct {
	service *service.Service
	repo    port.IRepo
}

func New(service *service.Service, repo port.IRepo) port.IUseCase {

	return &UseCase{
		service: service,
		repo:    repo}
}
