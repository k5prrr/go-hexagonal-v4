package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBConfig struct {
	Name     string
	User     string
	Password string
	Host     string
	Port     string
}

func (c *DBConfig) Validate() error {
	if c.Name == "" {
		return errors.New("database name is required")
	}
	if c.User == "" {
		return errors.New("database user is required")
	}
	if c.Password == "" {
		return errors.New("database password is required")
	}
	if c.Host == "" {
		return errors.New("database host is required")
	}
	if c.Port == "" {
		return errors.New("database port is required")
	}
	return nil
}

type IDB interface {
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, query string, args ...any) pgx.Row
	Close()
}

type pgxDB struct {
	pool *pgxpool.Pool
}

func New(conf *DBConfig) (IDB, error) {
	if conf == nil {
		return nil, errors.New("DBConfig is nil")
	}
	if err := conf.Validate(); err != nil {
		return nil, fmt.Errorf("invalid DB config: %w", err)
	}

	connectString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		conf.User, conf.Password, conf.Host, conf.Port, conf.Name,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, connectString)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &pgxDB{pool: pool}, nil
}

func (d *pgxDB) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	return d.pool.Query(ctx, query, args...)
}

func (d *pgxDB) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	return d.pool.Exec(ctx, query, args...)
}

func (d *pgxDB) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	return d.pool.QueryRow(ctx, query, args...)
}

func (d *pgxDB) Close() {
	if d.pool != nil {
		d.pool.Close()
	}
}
