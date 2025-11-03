package usecase

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
)

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

func (u *UseCase) AddChatIdByCodeCheckPhone(ctx context.Context, code string, chatId int) error {
	return nil
}

func (u *UseCase) AddPhoneByChatId(ctx context.Context, chatId int, phone string) error {
	return nil
}

func randStringHex(nBytes int) (string, error) {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
