package port

import "context"

// Так, давай на первых этапах IUseCase будет как общий, а потом разделим
type IUseCase interface {
	CreateCodeCheckPhone(ctx context.Context, t string) (string, error)
	AddUserIdCodeCheckPhone(ctx context.Context, code string, id int)
	AddPhoneByChatId(ctx context.Context, id int, s string)
}
type IRepo interface {
	AddCodeCheckPhone(ctx context.Context, code, codeType string) error
}
