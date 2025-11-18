package postgres

import (
	"app/internal/app/core/domain"
	"app/pkg/database"
	"app/pkg/utilities"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
)

type Repo struct {
	db database.IDB
}

func New(db database.IDB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) AddCodeCheckPhone(ctx context.Context, code, codeType string) error {
	query := `
		INSERT INTO auth_codes (code, type, uuid, phone, created_at, updated_at)
		VALUES ($1, $2, NULL, NULL, $3, $3)
	`

	_, err := r.db.Exec(ctx, query, code, codeType, time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) GetCodeByCode(ctx context.Context, code string) (*domain.AuthCode, error) {
	query := `
		SELECT id, code, type, uuid, phone, created_at, updated_at
		FROM auth_codes
		WHERE code = $1
	`
	row := r.db.QueryRow(ctx, query, code)
	var ac domain.AuthCode
	var uuidStr, phoneStr *string
	err := row.Scan(
		&ac.ID,
		&ac.Code,
		&ac.Type,
		&uuidStr,
		&phoneStr,
		&ac.CreatedAt,
		&ac.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("code not found")
		}
		return nil, fmt.Errorf("failed to get code: %w", err)
	}
	ac.UUID = uuidStr
	ac.Phone = phoneStr
	return &ac, nil
}

func (r *Repo) UpdateChatIDByCode(ctx context.Context, code string, chatID int) error {
	query := `
		UPDATE auth_codes
		SET uuid = $1, updated_at = $2
		WHERE code = $3
	`
	_, err := r.db.Exec(ctx, query, strconv.Itoa(chatID), time.Now(), code)
	if err != nil {
		return fmt.Errorf("failed to update chat_id: %w", err)
	}
	return nil
}

func (r *Repo) UpdatePhoneByChatID(ctx context.Context, chatID int, phone string) error {
	query := `
		UPDATE auth_codes
		SET phone = $1, updated_at = $2
		WHERE uuid = $3
	`
	_, err := r.db.Exec(ctx, query, phone, time.Now(), strconv.Itoa(chatID))
	if err != nil {
		return fmt.Errorf("failed to update phone: %w", err)
	}
	return nil
}

func (r *Repo) CreateUser(ctx context.Context, tgID int64, phone int64) (int64, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx) // откатится, если не вызван tx.Commit()

	// 1. Создаём пользователя (телефон — строка, а не int64!)
	var userID int64
	err = tx.QueryRow(ctx, `
		INSERT INTO users (phone, created_at, updated_at)
		VALUES ($1, NOW(), NOW())
		RETURNING id
	`, strconv.FormatInt(phone, 10)).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert user: %w", err)
	}


	tokenTime := r.db.Time() // for UNIQUE
	token := fmt.Sprintf(
		"%s%s",
		tokenTime,
		utilities.RandomString(64 - len(tokenTime)),
		)

	// 2. Создаём запись в auth
	_, err = tx.Exec(ctx,
		`
		INSERT INTO auth (user_id, tg_id, token, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		`,
		strconv.FormatInt(userID, 10),
		strconv.FormatInt(tgID, 10),
		token,
		)
	if err != nil {
		return 0, fmt.Errorf("failed to insert auth: %w", err)
	}

	// 3. Коммитим транзакцию
	if err = tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return userID, nil
}
