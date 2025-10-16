package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"app/internal/app/core/port"
)

type Auth struct {
	repo *port.IAuthRepo
}

func NewAuth(repo port.IAuthRepo) *Auth {
	return &Auth{repo: repo}
}

func (a *Auth) CreateCodeForCheckPhone(ctx context.Context, t string) (string, error) {
	if t != "registration" && t != "client" {
		return "", errors.New("invalid type")
	}

	code, err := randStringHex(16)
	if err != nil {
		return "", fmt.Errorf("generate code: %w", err)
	}

	if err := a.repo.AddCode(code); err != nil {
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
