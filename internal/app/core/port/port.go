package port

import "context"

// Так, давай на первых этапах IUseCase будет как общий, а потом разделим
type IUseCase interface {
	CreateCodeCheckPhone(ctx context.Context, t string) (string, error)
	AddChatIdByCodeCheckPhone(ctx context.Context, code string, id int) error
	AddPhoneByChatId(ctx context.Context, id int, s string) error
}
type IRepo interface {
	AddCodeCheckPhone(ctx context.Context, code, codeType string) error
}
