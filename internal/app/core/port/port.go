package port

import (
	"app/internal/app/core/domain"
	"context"
)

// Так, давай на первых этапах IUseCase будет как общий, а потом разделим
type IUseCase interface {
	CreateCodeCheckPhone(ctx context.Context, t string) (string, error)
	AddChatIdByCodeCheckPhone(ctx context.Context, code string, id int) error
	AddPhoneByChatId(ctx context.Context, id int, s string) error
	SendPasswordByPhone(ctx context.Context, phone int64) error
}
type IRepo interface {
	AddCodeCheckPhone(ctx context.Context, code, codeType string) error
	GetCodeByCode(ctx context.Context, code string) (*domain.AuthCode, error)
	UpdateChatIDByCode(ctx context.Context, code string, chatID int) error
	UpdatePhoneByChatID(ctx context.Context, chatID int, phone string) error
}
