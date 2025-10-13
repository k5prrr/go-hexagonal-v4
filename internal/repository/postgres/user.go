package postgres

import (
	"context"
	"fmt"
	"time"

	"app/internal/domain/user"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UserRepositoryPostgres struct {
	conn *pgx.Conn
}

func NewUserRepositoryPostgres(conn *pgx.Conn) *UserRepositoryPostgres {
	return &UserRepositoryPostgres{
		conn: conn,
	}
}

func (u *UserRepositoryPostgres) All(ctx context.Context) ([]*user.User, error) {
	if u.conn == nil {
		return nil, fmt.Errorf("database connection is not established")
	}

	rows, err := u.conn.Query(ctx, `
        SELECT 
            uid, family_name, name, middle_name, birth_date, phone, email, 
            phone_confirmed, email_confirmed, last_login, created_at, updated_at, deleted_at,
            description, password_hash, key_api
        FROM users
        WHERE deleted_at IS NULL
    `)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}
	defer rows.Close()

	var users []*user.User
	for rows.Next() {
		var userModel user.User
		err := rows.Scan(
			&userModel.Uid,
			&userModel.FamilyName,
			&userModel.Name,
			&userModel.MiddleName,
			&userModel.BirthDate,
			&userModel.Phone,
			&userModel.Email,
			&userModel.PhoneConfirmed,
			&userModel.EmailConfirmed,
			&userModel.LastLogin,
			&userModel.CreatedAt,
			&userModel.UpdatedAt,
			&userModel.DeletedAt,
			&userModel.Description,
			&userModel.PasswordHash,
			&userModel.KeyApi,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		users = append(users, &userModel)
	}

	return users, nil
}

func (u *UserRepositoryPostgres) Create(ctx context.Context, userModel *user.User) (string, error) {
	if u.conn == nil {
		return "", fmt.Errorf("database connection is not established")
	}

	userModel.Uid = uuid.New().String()
	now := time.Now()
	userModel.CreatedAt = now
	userModel.UpdatedAt = now

	_, err := u.conn.Exec(ctx, `
        INSERT INTO users (
            uid, family_name, name, middle_name, birth_date, phone, email, 
            phone_confirmed, email_confirmed, created_at, updated_at, description,
            password_hash, key_api
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
    `,
		userModel.Uid,
		userModel.FamilyName,
		userModel.Name,
		userModel.MiddleName,
		userModel.BirthDate,
		userModel.Phone,
		userModel.Email,
		userModel.PhoneConfirmed,
		userModel.EmailConfirmed,
		userModel.CreatedAt,
		userModel.UpdatedAt,
		userModel.Description,
		userModel.PasswordHash,
		userModel.KeyApi,
	)

	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	return userModel.Uid, nil
}

func (u *UserRepositoryPostgres) Read(ctx context.Context, uid string) (*user.User, error) {
	if u.conn == nil {
		return nil, fmt.Errorf("database connection is not established")
	}

	var userModel user.User

	err := u.conn.QueryRow(ctx, `
        SELECT 
            uid, family_name, name, middle_name, birth_date, phone, email, 
            phone_confirmed, email_confirmed, last_login, created_at, updated_at, deleted_at,
            description, password_hash, key_api
        FROM users
        WHERE uid = $1 AND deleted_at IS NULL
    `, uid).Scan(
		&userModel.Uid,
		&userModel.FamilyName,
		&userModel.Name,
		&userModel.MiddleName,
		&userModel.BirthDate,
		&userModel.Phone,
		&userModel.Email,
		&userModel.PhoneConfirmed,
		&userModel.EmailConfirmed,
		&userModel.LastLogin,
		&userModel.CreatedAt,
		&userModel.UpdatedAt,
		&userModel.DeletedAt,
		&userModel.Description,
		&userModel.PasswordHash,
		&userModel.KeyApi,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to read user: %w", err)
	}

	return &userModel, nil
}

func (u *UserRepositoryPostgres) Update(ctx context.Context, userModel *user.User) error {
	if u.conn == nil {
		return fmt.Errorf("database connection is not established")
	}

	if userModel.Uid == "" {
		return fmt.Errorf("user UID is required for update")
	}

	userModel.UpdatedAt = time.Now()

	query := `
        UPDATE users SET
            family_name = $1,
            name = $2,
            middle_name = $3,
            birth_date = $4,
            phone = $5,
            email = $6,
            phone_confirmed = $7,
            email_confirmed = $8,
            updated_at = $9,
            description = $10,
            password_hash = $11,
            key_api = $12
        WHERE uid = $13 AND deleted_at IS NULL
        RETURNING uid`

	var existingUID string
	err := u.conn.QueryRow(ctx, query,
		userModel.FamilyName,
		userModel.Name,
		userModel.MiddleName,
		userModel.BirthDate,
		userModel.Phone,
		userModel.Email,
		userModel.PhoneConfirmed,
		userModel.EmailConfirmed,
		userModel.UpdatedAt,
		userModel.Description,
		userModel.PasswordHash,
		userModel.KeyApi,
		userModel.Uid,
	).Scan(&existingUID)

	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("user with UID %s not found or already deleted", userModel.Uid)
		}
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (u *UserRepositoryPostgres) Delete(ctx context.Context, uid string) error {
	if u.conn == nil {
		return fmt.Errorf("database connection is not established")
	}

	deletedAt := time.Now()

	cmdTag, err := u.conn.Exec(ctx, `
        UPDATE users SET
            deleted_at = $1
        WHERE uid = $2 AND deleted_at IS NULL
    `, deletedAt, uid)

	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no user was deleted (maybe already deleted or not exists)")
	}

	return nil
}
