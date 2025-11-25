package postgres

/*
	1 Заменяем все RepoUser
	2 Заменяем все domain.User
	3 Меняем поля в Change (Важно в порядке columns)

	В таблице всегда должны быть
	id, created_at, updated_at, deleted_at
*/
import (
	"app/internal/app/core/domain"
	"app/pkg/database"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrUnsupported  = errors.New("unsupported filter key")
)

type RepoUser struct {
	db               database.IDB
	tableName        string
	columns          []string
	columnsStr       string
	columnsStrUpdate string
	columnsLen       int
	allowedFilters   map[string]struct{}
	softDelete	bool
}

// Change
func NewRepoUser(db database.IDB) *RepoUser {
	return &RepoUser{
		db:        db,
		tableName: "users",
		columns: []string{"family_name", "name", "middle_name", "phone", "email",
			"birth_date", "parent_id", "gender_id", "role_id",
		},
		allowedFilters: map[string]struct{}{
			"id":          {},
			"phone":       {},
			"email":       {},
			"family_name": {},
			"name":        {},
		},
		softDelete: true,
	}
}
func (r *RepoUser) scanEntityRow(row pgx.Row) (*domain.User, error) {
	var e domain.User

	err := row.Scan(
		&e.ID,

		&e.FamilyName,
		&e.Name,
		&e.MiddleName,
		&e.Phone,
		&e.Email,
		&e.BirthDate,
		&e.ParentID,
		&e.GenderID,
		&e.RoleID,

		&e.CreatedAt,
		&e.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("scan entity row: %w", err)
	}

	return &e, nil
}
func (r *RepoUser) Add(ctx context.Context, entity *domain.User) (int64, error) {
	if entity == nil {
		return 0, fmt.Errorf("%w: entity is nil", ErrInvalidInput)
	}

	now := time.Now()
	entity.CreatedAt = now
	entity.UpdatedAt = now

	query := fmt.Sprintf(`
		INSERT INTO %s (%s) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id
	`, r.tableName, r.getColumnsStr())

	var id int64
	err := r.db.QueryRow(ctx, query,

		entity.FamilyName,
		entity.Name,
		entity.MiddleName,
		entity.Phone,
		entity.Email,
		entity.BirthDate,
		entity.ParentID,
		entity.GenderID,
		entity.RoleID,

		entity.CreatedAt,
		entity.UpdatedAt,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("insert entity: %w", err)
	}

	entity.ID = id

	return id, nil
}
func (r *RepoUser) Update(ctx context.Context, id int64, entity *domain.User) error {
	if id <= 0 {
		return fmt.Errorf("%w: invalid id %d", ErrInvalidInput, id)
	}
	if entity == nil {
		return fmt.Errorf("%w: entity is nil", ErrInvalidInput)
	}

	entity.UpdatedAt = time.Now()

	query := fmt.Sprintf(`
		UPDATE %s SET %s
		WHERE id = $%d AND deleted_at IS NULL
	`, r.tableName, r.getColumnsStrUpdate(), r.getColumnsUpdateI())

	_, err := r.db.Exec(ctx, query,

		entity.FamilyName,
		entity.Name,
		entity.MiddleName,
		entity.Phone,
		entity.Email,
		entity.BirthDate,
		entity.ParentID,
		entity.GenderID,
		entity.RoleID,

		entity.UpdatedAt,
		id,
	)

	if err != nil {
		return fmt.Errorf("update entity %d: %w", id, err)
	}

	return nil
}
func (r *RepoUser) UpdateBy(ctx context.Context, filterKey, filterValue string, entity *domain.User) error {
	if entity == nil {
		return fmt.Errorf("%w: entity is nil", ErrInvalidInput)
	}
	if err := r.validateFilterKey(filterKey); err != nil {
		return err
	}

	entity.UpdatedAt = time.Now()

	query := fmt.Sprintf(`
		UPDATE %s SET %s
		WHERE deleted_at IS NULL AND %s = $%d
	`, r.tableName, r.getColumnsStrUpdate(), filterKey, r.getColumnsUpdateI())

	_, err := r.db.Exec(ctx, query,

		entity.FamilyName,
		entity.Name,
		entity.MiddleName,
		entity.Phone,
		entity.Email,
		entity.BirthDate,
		entity.ParentID,
		entity.GenderID,
		entity.RoleID,

		entity.UpdatedAt,
		filterValue,
	)

	if err != nil {
		return fmt.Errorf("update entity by %s=%s: %w", filterKey, filterValue, err)
	}

	return nil
}

// Not Change
func (r *RepoUser) Get(ctx context.Context, id int64) (*domain.User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: invalid id %d", ErrInvalidInput, id)
	}

	query := fmt.Sprintf(`
		SELECT id, %s
		FROM %s
		WHERE id = $1 AND deleted_at IS NULL
		LIMIT 1
	`, r.getColumnsStr(), r.tableName)

	row := r.db.QueryRow(ctx, query, id)

	return r.scanEntityRow(row)
}
func (r *RepoUser) GetBy(ctx context.Context, filterKey, filterValue string) (*domain.User, error) {
	if err := r.validateFilterKey(filterKey); err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`
		SELECT id, %s
		FROM %s
		WHERE deleted_at IS NULL AND %s = $1
		LIMIT 1
	`, r.getColumnsStr(), r.tableName, filterKey)

	row := r.db.QueryRow(ctx, query, filterValue)

	return r.scanEntityRow(row)
}
func (r *RepoUser) List(ctx context.Context, offset, limit int64) ([]domain.User, error) {
	var queryEnd string
	if offset != 0 || limit != 0 {
		queryEnd = fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
	}
	query := fmt.Sprintf(`
		SELECT id, %s
		FROM %s
		WHERE deleted_at IS NULL
		%s
	`, r.getColumnsStr(), r.tableName, queryEnd)

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query entities: %w", err)
	}

	return r.scanEntityRows(rows)
}
func (r *RepoUser) ListBy(ctx context.Context, filterKey, filterValue string, offset, limit int64) ([]domain.User, error) {
	var queryEnd string
	if offset != 0 || limit != 0 {
		queryEnd = fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
	}

	if err := r.validateFilterKey(filterKey); err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`
		SELECT %s
		FROM %s
		WHERE deleted_at IS NULL AND %s = $1
		%s
	`, r.getColumnsStr(), r.tableName, filterKey, queryEnd)

	rows, err := r.db.Query(ctx, query, filterValue)
	if err != nil {
		return nil, fmt.Errorf("query entities by %s: %w", filterKey, err)
	}

	return r.scanEntityRows(rows)
}
func (r *RepoUser) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("%w: invalid id %d", ErrInvalidInput, id)
	}

	query := fmt.Sprintf(`
		UPDATE %s SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`, r.tableName)

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("soft-delete entity %d: %w", id, err)
	}

	return nil
}
func (r *RepoUser) DeleteBy(ctx context.Context, filterKey, filterValue string) error {
	if err := r.validateFilterKey(filterKey); err != nil {
		return err
	}

	query := fmt.Sprintf(`
		UPDATE %s SET deleted_at = NOW()
		WHERE deleted_at IS NULL AND %s = $1
	`, r.tableName, filterKey)

	_, err := r.db.Exec(ctx, query, filterValue)
	if err != nil {
		return fmt.Errorf("soft-delete entity by %s=%s: %w", filterKey, filterValue, err)
	}

	return nil
}

// Not Change Utils
func (r *RepoUser) getColumnsStr() string {
	if r.columnsStr == "" {
		r.columnsStr = fmt.Sprintf(
			"%s, created_at, updated_at",
			strings.Join(r.columns, ", "),
		)
	}
	return r.columnsStr
}
func (r *RepoUser) getColumnsStrUpdate() string {
	if r.columnsStrUpdate == "" {
		tmpI := len(r.columns)
		tmp := make([]string, tmpI)
		for i, name := range r.columns {
			tmp[i] = fmt.Sprintf("%s = $%d", name, i+1)
		}
		r.columnsStrUpdate = fmt.Sprintf(
			"%s, updated_at = $%d",
			strings.Join(tmp, ", "),
			tmpI+1,
		)
	}
	return r.columnsStrUpdate
}
func (r *RepoUser) getColumnsUpdateI() int {
	if r.columnsLen == 0 {
		r.columnsLen = len(r.columns) + 2
	}

	return r.columnsLen
}
func (r *RepoUser) validateFilterKey(key string) error {
	if _, ok := r.allowedFilters[key]; !ok {
		return fmt.Errorf("%w: %q", ErrUnsupported, key)
	}
	return nil
}
func (r *RepoUser) scanEntityRows(rows pgx.Rows) ([]domain.User, error) {
	defer rows.Close()

	var entities []domain.User
	for rows.Next() {
		e, err := r.scanEntityRow(rows)
		if err != nil {
			return nil, fmt.Errorf("scan entity in list: %w", err)
		}
		entities = append(entities, *e)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return entities, nil
}
