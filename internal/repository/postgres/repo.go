package postgres

import (
	"app/pkg/database"
	"context"
	"time"
)

type Repo struct {
	db database.IDB
}

func NewRepo(db database.IDB) *Repo {
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
