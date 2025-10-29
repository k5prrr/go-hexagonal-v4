package usecase

import (
	"app/internal/app/core/port"
	"encoding/hex"
	"fmt"

	"context"
	"crypto/rand"

	"errors"
)

// Тут он как основной, тут только создания
type UseCase struct {
	repo port.IRepo
}

func NewUseCase(repo port.IRepo) port.IUseCase {
	return &UseCase{repo: repo}
}

func (u *UseCase) CreateCodeCheckPhone(ctx context.Context, action string) (string, error) {
	if action != "registration" && action != "client" {
		return "", errors.New("invalid action")
	}

	code, err := randStringHex(16)
	if err != nil {
		return "", fmt.Errorf("generate code: %w", err)
	}

	if err := u.repo.AddCodeCheckPhone(ctx, code, action); err != nil {
		return "", fmt.Errorf("save code: %w", err)
	}

	return code, nil
}

func randStringHex(nBytes int) (string, error) {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
