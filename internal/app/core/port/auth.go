package port

import "context"

type IAuth interface {
	CreateCodeForCheckPhone(ctx context.Context, t string) (string, error)
}
type IAuthRepo interface {
	AddCode(ctx context.Context, code, codeType string) error
}
