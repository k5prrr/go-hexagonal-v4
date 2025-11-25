package port

import (
	"app/internal/app/core/domain"
	"context"
)

// Так, давай на первых этапах IUseCase будет как общий, а потом разделим
type IUseCase interface {
	/*CreateCodeCheckPhone(ctx context.Context, t string) (string, error)
	AddChatIdByCodeCheckPhone(ctx context.Context, code string, id int) error
	AddPhoneByChatId(ctx context.Context, id int, s string) error*/

	//SendAuthCode(ctx context.Context, phone int64) error
	CreateUser(ctx context.Context, tgID, phone int64) (int64, error)
	/*LoginCode(ctx context.Context, phone int64, code int64) (int64, string, error)

	 */
	Test1() (*domain.User, error)

}
type IRepo interface {
	/*AddCodeCheckPhone(ctx context.Context, code, codeType string) error
	GetCodeByCode(ctx context.Context, code string) (*domain.AuthCode, error)
	UpdateChatIDByCode(ctx context.Context, code string, chatID int) error
	UpdatePhoneByChatID(ctx context.Context, chatID int, phone string) error
	*/

	/*
		CreateUser(ctx context.Context, tgID int64, phone int64) (int64, error)
		UserByPhone(ctx context.Context, phone int64) (*domain.User, error)
		AuthByUserID(ctx context.Context, userID int64) (*domain.Auth, error)
		UpdateAuthCode(id int64, code int) error
	*/

}

type Itg interface {
	SendPhoto(chatID int64, urlPhoto string, message string, replyMarkup string) (string, error)
	SendMessage(chatID int64, message string, replyMarkup string) (string, error)
}

type IRepoUser interface {
	Add(ctx context.Context, entity *domain.User) (int64, error)

	Get(ctx context.Context, id int64) (*domain.User, error)
	GetBy(ctx context.Context, filterKey, filterValue string) (*domain.User, error)

	List(ctx context.Context, offset, limit int64) ([]domain.User, error)
	ListBy(ctx context.Context, filterKey, filterValue string, offset, limit int64) ([]domain.User, error)

	Update(ctx context.Context, id int64, entity *domain.User) error
	UpdateBy(ctx context.Context, filterKey, filterValue string, entity *domain.User) error

	Delete(ctx context.Context, id int64) error
	DeleteBy(ctx context.Context, filterKey, filterValue string) error
}
