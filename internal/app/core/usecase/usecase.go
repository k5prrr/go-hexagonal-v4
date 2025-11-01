package usecase

import (
	"app/internal/app/core/port"
)

// Тут он как основной, тут только создания
type UseCase struct {
	repo port.IRepo
}

func NewUseCase(repo port.IRepo) port.IUseCase {
	return &UseCase{repo: repo}
}
