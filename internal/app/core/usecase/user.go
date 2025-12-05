package usecase

import (
	"app/internal/app/core/domain"
	"context"
	"errors"
	"fmt"
	"time"
)

/*
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

	//	Добавляет тому, кто есть,
	//	если нету, то возврат ошибки

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
*/

/*
	func randStringHex(nBytes int) (string, error) {
		b := make([]byte, nBytes)
		if _, err := rand.Read(b); err != nil {
			return "", err
		}
		return hex.EncodeToString(b), nil
	}

	func (u *UseCase) userAuthByPhone(ctx context.Context, phone int64) (*domain.User, *domain.Auth, error) {
		user, err := u.repo.UserByPhone(ctx, phone)
		if err != nil {
			return nil, nil, fmt.Errorf("get user by phone: %w", err)
		}

		auth, err := u.repo.AuthByUserID(ctx, user.ID)
		if err != nil {
			return nil, nil, fmt.Errorf("get auth by user id: %w", err)
		}

		return user, auth, nil
	}

func (u *UseCase) SendAuthCode(ctx context.Context, phone int64) error {

		_, auth, err := u.userAuthByPhone(ctx, phone)
		if err != nil {
			return fmt.Errorf("get user auth by phone: %w", err)
		}

		code := rand.Intn(900000) + 100000

		if err = u.repo.UpdateAuthCode(auth.ID, code); err != nil {
			return fmt.Errorf("update auth code: %w", err)
		}

		if _, err = u.tg.SendMessage(
			auth.TgID,
			fmt.Sprintf("Код: %d", code),
			"",
			); err != nil {
			return fmt.Errorf("send message: %w", err)
		}

		return nil
	}
*/
func (u *UseCase) CreateUser(ctx context.Context, tgID int64, phone string) (int64, error) {
	var err error

	// Есть ли телефон, то вылет
	user, _ := u.service.RepoUser.GetBy(ctx, "phone", phone)
	if user != nil {
		return 0, errors.New("user have by phone")
	}

	// Есть ли тг, то обновим телефон
	auth, _ := u.service.RepoAuth.GetByInt(ctx, "tg_id", tgID)
	if auth != nil {
		err = u.service.RepoUser.UpdateColumn(ctx, auth.UserID, "phone", phone)
		if err != nil {
			return 0, errors.New("user not update phone")
		}

		return auth.UserID, nil
	}

	// создаём секрет
	secret, err := u.service.Secret(32) // 32 to 64 chars
	if err != nil {
		return 0, err
	}

	userID, err := u.service.RepoUser.Add(ctx, &domain.User{Phone: phone})
	if err != nil {
		return 0, errors.New("user not add phone")
	}

	_, err = u.service.RepoAuth.Add(ctx, &domain.Auth{
		UserID: userID,
		TgID:   tgID,
		Secret: secret,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to add auth for tg: %w", err)
	}

	return userID, nil
}

// Отправка кода для входа
func (u *UseCase) SendAuthCode(ctx context.Context, phone string) error {
	userFull, err := u.service.UserAuthByPhone(ctx, phone)
	if err != nil {
		return fmt.Errorf("err UserAuthByPhone in SendAuthCode: %w", err)
	}
	auth := userFull.Auth

	if time.Since(auth.UpdatedAt) < 6*time.Second {
		return fmt.Errorf("err time load")
	}

	code := u.service.Code()

	err = u.service.RepoAuth.UpdateColumn(ctx, auth.ID, "code", code)
	if err != nil {
		return fmt.Errorf("err UpdateColumn in SendAuthCode: %w", err)
	}

	_, err = u.service.Tg.SendMessage(auth.TgID, fmt.Sprintf("Code: %s", code), "")
	if err != nil {
		return fmt.Errorf("err SendMessage in SendAuthCode: %w", err)
	}

	return nil
}
func (u *UseCase) CheckAuthCode(ctx context.Context, phone, code string) (string, error) {
	userFull, err := u.service.UserAuthByPhone(ctx, phone)
	if err != nil {
		return "", fmt.Errorf("err UserAuthByPhone in SendAuthCode: %w", err)
	}
	auth := userFull.Auth

	if auth.Code == nil {
		return "", fmt.Errorf("not code")
	}

	err = u.service.RepoAuth.UpdateColumn(ctx, auth.ID, "code", "NULL")
	if err != nil {
		return "", fmt.Errorf("err UpdateColumn in SendAuthCode: %w", err)
	}

	if *auth.Code != code {
		return "", fmt.Errorf("not code")
	}

	token := u.service.Token(auth.ID, auth.Secret)

	return token, nil
}

func (u *UseCase) CurrentUser(ctx context.Context, token string) (*domain.UserFull, error) {
	return u.service.CurrentUser(ctx, token)
}
