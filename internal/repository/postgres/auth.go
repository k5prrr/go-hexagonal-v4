package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/jackc/pgx/v5"
)

type RepoAuth struct {
	db *pgx.Conn
}

func NewRepoAuth(conn *pgx.Conn) *RepoAuth {
	return &RepoAuth{db: conn}
}

func (r *RepoAuth) AddCode(ctx context.Context, code, codeType string) error {
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
