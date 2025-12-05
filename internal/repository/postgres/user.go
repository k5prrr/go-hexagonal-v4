package postgres

/*
	1 Заменяем все RepoUser
	2 Заменяем все domain.User
	3 Меняем поля в Change (Важно в порядке columns)

	В таблице всегда должны быть
	id, created_at, updated_at, deleted_at
	но в columns их НЕ указываем
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
type RepoUser struct {
	db               database.IDB
	tableName        string
	columns          []string
	columnsStr       string
	columnsStrUpdate string
	columnsInput     string
	columnsUpdateI   int
	columnsMap       map[string]struct{}
}

// Change
func NewRepoUser(db database.IDB) *RepoUser {
	return &RepoUser{
		db:        db,
		tableName: "users",
		columns: []string{
			"family_name", "name", "middle_name", "phone", "email",
			"birth_date", "parent_id", "gender_id", "role_id",
		},
	}
}
func (r *RepoUser) scanEntityRow(row pgx.Row) (*domain.User, error) {
	var entity domain.User

	err := row.Scan(
		&entity.ID,

		&entity.FamilyName,
		&entity.Name,
		&entity.MiddleName,
		&entity.Phone,
		&entity.Email,
		&entity.BirthDate,
		&entity.ParentID,
		&entity.GenderID,
		&entity.RoleID,

		&entity.CreatedAt,
		&entity.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("scan entity row: %w", err)
	}

	return &entity, nil
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
		VALUES (%s)
		RETURNING id
	`, r.tableName, r.getColumnsStr(), r.getColumnsInput())

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
	if err := r.validateColumn(filterKey); err != nil {
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
	if err := r.validateColumn(filterKey); err != nil {
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
func (r *RepoUser) GetByInt(ctx context.Context, filterKey string, filterValue int64) (*domain.User, error) {
	if err := r.validateColumn(filterKey); err != nil {
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
		queryEnd = fmt.Sprintf("ORDER BY id LIMIT %d OFFSET %d", limit, offset)
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
		queryEnd = fmt.Sprintf("ORDER BY id LIMIT %d OFFSET %d", limit, offset)
	}

	if err := r.validateColumn(filterKey); err != nil {
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
func (r *RepoUser) UpdateColumn(ctx context.Context, id int64, key, value string) error {
	var err error
	if err = r.validateColumn(key); err != nil {
		return err
	}

	var query string
	if value == "NULL" {
		query = fmt.Sprintf(`
			UPDATE %s SET %s = NULL, updated_at = $1
			WHERE id = $2 AND deleted_at IS NULL
		`, r.tableName, key)
		_, err = r.db.Exec(ctx, query,
			time.Now(),
			id,
		)
	} else {
		query = fmt.Sprintf(`
			UPDATE %s SET %s = $1, updated_at = $2
			WHERE id = $3 AND deleted_at IS NULL
		`, r.tableName, key)
		_, err = r.db.Exec(ctx, query,
			value,
			time.Now(),
			id,
		)
	}

	if err != nil {
		return fmt.Errorf("UpdateColumn %d %s=%s: %w", id, key, value, err)
	}

	return nil
}
func (r *RepoUser) Delete(ctx context.Context, id int64, soft bool) error {
	if id <= 0 {
		return fmt.Errorf("%w: invalid id %d", ErrInvalidInput, id)
	}

	var query string
	if soft {
		query = "UPDATE %s SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL"
	} else {
		query = "DELETE FROM %s WHERE id = $1"
	}
	query = fmt.Sprintf(query, r.tableName)

	res, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete entity %d: %w", id, err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("%w: %d not found or already deleted", ErrNotFound, id)
	}

	return nil
}
func (r *RepoUser) DeleteBy(ctx context.Context, filterKey, filterValue string, soft bool) error {
	if err := r.validateColumn(filterKey); err != nil {
		return err
	}

	var query string
	if soft {
		query = "UPDATE %s SET deleted_at = NOW() WHERE deleted_at IS NULL AND %s = $1"
	} else {
		query = "DELETE FROM %s WHERE deleted_at IS NULL AND %s = $1"
	}
	query = fmt.Sprintf(query, r.tableName, filterKey)

	res, err := r.db.Exec(ctx, query, filterValue)
	if err != nil {
		return fmt.Errorf("delete entity by %s=%s: %w", filterKey, filterValue, err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("%w: not found or already deleted", ErrNotFound)
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
	if r.columnsUpdateI == 0 {
		r.columnsUpdateI = len(r.columns) + 2
	}

	return r.columnsUpdateI
}
func (r *RepoUser) getColumnsInput() string {
	if r.columnsInput == "" {
		var result []string
		kol := r.getColumnsUpdateI()
		for i := 0; i < kol; i++ {
			result = append(result, fmt.Sprintf("$%d", i+1))
		}
		r.columnsInput = strings.Join(result, ", ")
	}

	return r.columnsInput
}
func (r *RepoUser) getColumnsMap() map[string]struct{} {
	if r.columnsMap == nil {
		r.columnsMap = make(map[string]struct{})
		for _, name := range r.columns {
			r.columnsMap[name] = struct{}{}
		}
	}

	return r.columnsMap
}
func (r *RepoUser) validateColumn(key string) error {
	columnsMap := r.getColumnsMap()
	if _, ok := columnsMap[key]; !ok {
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
