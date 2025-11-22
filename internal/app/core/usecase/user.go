package usecase

import (
	"context"
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
func (u *UseCase) CreateUser(ctx context.Context, tgID, phone int64) (int64, error) {
	return 0, nil
	/*_, _, _ = ctx, tgID, phone
	id, err := u.repo.CreateUser(ctx, tgID, phone)
	if err != nil {

		return 0, fmt.Errorf("failed repo CreateUser: %w", err)
	}

	return id, nil

	 */
}

/*
func (u *UseCase) LoginCode(ctx context.Context, phone int64, code int64) (int64, string, error) {
	return 0, "", nil
}
*/

