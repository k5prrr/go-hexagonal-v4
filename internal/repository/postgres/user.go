package postgres
/*
1 Меняем все RepoUser
2 Меняем все domain.User
3 Меняем поля в Change

*/
import (
	"app/internal/app/core/domain"
	"app/pkg/database"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrUnsupported  = errors.New("unsupported filter key")
)

// Change
var allowedFilters = map[string]struct{}{
	"id":          {},
	"phone":       {},
	"email":       {},
	"family_name": {},
	"name":        {},
}

func validateFilterKey(key string) error {
	if _, ok := allowedFilters[key]; !ok {
		return fmt.Errorf("%w: %q", ErrUnsupported, key)
	}
	return nil
}

type RepoUser struct {
	db        database.IDB
	tableName string
	columns   []string
}

// Change
func NewRepoUser(db database.IDB) *RepoUser {
	return &RepoUser{
		db:        db,
		tableName: "users",
		columns:   []string{"family_name", "name", "middle_name", "phone", "email",
			"birth_date", "parent_id", "gender_id", "role_id",
			"created_at", "updated_at"},
	}
}

// Change
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

// Change
func (r *RepoUser) Add(ctx context.Context, entity *domain.User) (int64, error) {
	if entity == nil {
		return 0, fmt.Errorf("%w: entity is nil", ErrInvalidInput)
	}

	now := time.Now()
	entity.CreatedAt = now
	entity.UpdatedAt = now

	query := fmt.Sprintf(`
		INSERT INTO %s (
			family_name, name, middle_name, phone, email,
			birth_date, parent_id, gender_id, role_id,
			created_at, updated_aсt
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id
	`, r.tableName)

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

// !!!+решим
func (r *RepoUser) Get(ctx context.Context, id int64) (*domain.User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: invalid id %d", ErrInvalidInput, id)
	}

	query := fmt.Sprintf(`
		SELECT id, family_name, name, middle_name, phone,
		       email, birth_date, parent_id, gender_id, role_id,
		       created_at, updated_at
		FROM %s
		WHERE id = $1 AND deleted_at IS NULL
	`, r.tableName)

	row := r.db.QueryRow(ctx, query, id)

	return r.scanEntityRow(row)
}
func (r *RepoUser) GetBy(ctx context.Context, filterKey, filterValue string) (*domain.User, error) {
	if err := validateFilterKey(filterKey); err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`
		SELECT id, family_name, name, middle_name, phone,
		       email, birth_date, parent_id, gender_id, role_id,
		       created_at, updated_at
		FROM %s
		WHERE deleted_at IS NULL AND %s = $1
	`, r.tableName, filterKey)

	row := r.db.QueryRow(ctx, query, filterValue)

	return r.scanEntityRow(row)
}

// !!!+Решим
func (r *RepoUser) List(ctx context.Context, offset, limit int64) ([]domain.User, error) {
	var queryEnd string
	if offset != 0 || limit != 0 {
		queryEnd = fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
	}
	query := fmt.Sprintf(`
		SELECT id, family_name, name, middle_name, phone,
		       email, birth_date, parent_id, gender_id, role_id,
		       created_at, updated_at
		FROM %s
		WHERE deleted_at IS NULL
		%s
	`, r.tableName, queryEnd)

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

	if err := validateFilterKey(filterKey); err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`
		SELECT id, family_name, name, middle_name, phone,
		       email, birth_date, parent_id, gender_id, role_id,
		       created_at, updated_at
		FROM %s
		WHERE deleted_at IS NULL AND %s = $1
		%s
	`, r.tableName, filterKey, queryEnd)

	rows, err := r.db.Query(ctx, query, filterValue)
	if err != nil {
		return nil, fmt.Errorf("query entities by %s: %w", filterKey, err)
	}

	return r.scanEntityRows(rows)
}

// Change
func (r *RepoUser) Update(ctx context.Context, id int64, entity *domain.User) error {
	if id <= 0 {
		return fmt.Errorf("%w: invalid id %d", ErrInvalidInput, id)
	}
	if entity == nil {
		return fmt.Errorf("%w: entity is nil", ErrInvalidInput)
	}

	entity.UpdatedAt = time.Now()

	query := fmt.Sprintf(`
		UPDATE %s SET
		    family_name = $1,
		    name = $2,
		    middle_name = $3,
		    phone = $4,
		    email = $5,
		    birth_date = $6,
		    parent_id = $7,
		    gender_id = $8,
		    role_id = $9,
		    updated_at = $10
		WHERE id = $11 AND deleted_at IS NULL
	`, r.tableName)

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
func (r *RepoUser) UpdateBy(ctx context.Context, filterKey, filterValue string, entity *domain.User, limit int64) error {
	if entity == nil {
		return fmt.Errorf("%w: entity is nil", ErrInvalidInput)
	}
	if err := validateFilterKey(filterKey); err != nil {
		return err
	}
	var queryEnd string
	if limit != 0 {
		queryEnd = fmt.Sprintf("LIMIT %d", limit)
	}

	entity.UpdatedAt = time.Now()

	query := fmt.Sprintf(`
		UPDATE %s SET
		    family_name = $1,
		    name = $2,
		    middle_name = $3,
		    phone = $4,
		    email = $5,
		    birth_date = $6,
		    parent_id = $7,
		    gender_id = $8,
		    role_id = $9,
		    updated_at = $10
		WHERE deleted_at IS NULL AND %s = $11
		%s
	`, r.tableName, filterKey, queryEnd)

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
func (r *RepoUser) DeleteBy(ctx context.Context, filterKey, filterValue string, limit int64) error {
	if err := validateFilterKey(filterKey); err != nil {
		return err
	}
	var queryEnd string
	if limit != 0 {
		queryEnd = fmt.Sprintf("LIMIT %d", limit)
	}

	query := fmt.Sprintf(`
		UPDATE %s SET deleted_at = NOW()
		WHERE deleted_at IS NULL AND %s = $1
		%s
	`, r.tableName, filterKey, queryEnd)

	_, err := r.db.Exec(ctx, query, filterValue)
	if err != nil {
		return fmt.Errorf("soft-delete entity by %s=%s: %w", filterKey, filterValue, err)
	}

	return nil
}
