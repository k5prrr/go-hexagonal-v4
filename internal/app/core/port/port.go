package port

import (
	"app/internal/app/core/domain"
	"context"
)

// Так, давай на первых этапах IUseCase будет как общий, а потом разделим
type IUseCase interface {
	CreateUser(ctx context.Context, tgID int64, phone string) (int64, error)
	SendAuthCode(ctx context.Context, phone string) error
	CheckAuthCode(ctx context.Context, phone, code string) (string, error)

	CurrentUser(ctx context.Context, token string) (*domain.UserFull, error)

	Test1() (*domain.User, error)
}

type Itg interface {
	SendPhoto(chatID int64, urlPhoto string, message string, replyMarkup string) (string, error)
	SendMessage(chatID int64, message string, replyMarkup string) (string, error)
}

type IRepo interface{}
type IRepoUser interface {
	Add(ctx context.Context, entity *domain.User) (int64, error)

	Get(ctx context.Context, id int64) (*domain.User, error)
	GetBy(ctx context.Context, filterKey, filterValue string) (*domain.User, error)
	GetByInt(ctx context.Context, filterKey string, filterValue int64) (*domain.User, error)

	List(ctx context.Context, offset, limit int64) ([]domain.User, error)
	ListBy(ctx context.Context, filterKey, filterValue string, offset, limit int64) ([]domain.User, error)

	Update(ctx context.Context, id int64, entity *domain.User) error
	UpdateBy(ctx context.Context, filterKey, filterValue string, entity *domain.User) error
	UpdateColumn(ctx context.Context, id int64, key, value string) error

	Delete(ctx context.Context, id int64, soft bool) error
	DeleteBy(ctx context.Context, filterKey, filterValue string, soft bool) error
}
type IRepoAuth interface {
	Add(ctx context.Context, entity *domain.Auth) (int64, error)

	Get(ctx context.Context, id int64) (*domain.Auth, error)
	GetBy(ctx context.Context, filterKey, filterValue string) (*domain.Auth, error)
	GetByInt(ctx context.Context, filterKey string, filterValue int64) (*domain.Auth, error)

	List(ctx context.Context, offset, limit int64) ([]domain.Auth, error)
	ListBy(ctx context.Context, filterKey, filterValue string, offset, limit int64) ([]domain.Auth, error)

	Update(ctx context.Context, id int64, entity *domain.Auth) error
	UpdateBy(ctx context.Context, filterKey, filterValue string, entity *domain.Auth) error
	UpdateColumn(ctx context.Context, id int64, key, value string) error

	Delete(ctx context.Context, id int64, soft bool) error
	DeleteBy(ctx context.Context, filterKey, filterValue string, soft bool) error
}

type IMapAuth interface {
	Set(token string, userFull *domain.UserFull)
	Get(token string) *domain.UserFull
	Delete(token string)
}
