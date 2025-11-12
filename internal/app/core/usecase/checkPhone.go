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
	/*
		Добавляет тому, кто есть,
		если нету, то возврат ошибки
	*/
	_, err := u.repo.GetCodeByCode(ctx, code)
	if err != nil {
		return fmt.Errorf("code not found: %w", err)
	}

	// Обновляем UUID (в котором храним chatId как строку)
	err = u.repo.UpdateChatIDByCode(ctx, code, chatId)
	if err != nil {
		return fmt.Errorf("failed to update chat_id: %w", err)
	}

	return nil
}

func (u *UseCase) AddPhoneByChatId(ctx context.Context, chatId int, phone string) error {
	// Находим запись по chatId (хранится в uuid)
	err := u.repo.UpdatePhoneByChatID(ctx, chatId, phone)
	if err != nil {
		return fmt.Errorf("failed to update phone by chat_id: %w", err)
	}

	return nil
}

func randStringHex(nBytes int) (string, error) {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (*UseCase) SendPasswordByPhone(ctx context.Context, phone int64) error {
	return nil
}
