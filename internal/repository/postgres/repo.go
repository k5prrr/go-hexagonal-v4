package postgres

import (
	"app/internal/app/core/domain"
	"app/pkg/database"
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
