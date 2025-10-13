package database

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5"
)

type PostgresConfig struct {
	Name     string
	User     string
	Password string
	Host     string
	Port     string
}

func NewPostgresConnection(cfg *PostgresConfig) (*pgx.Conn, error) {
	if cfg.Host == "" {
		cfg.Host = "localhost"
	}
	if cfg.Port == "" {
		cfg.Port = "5432"
	}

	connectString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User,
		url.QueryEscape(cfg.Password),
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, connectString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return conn, nil
}

func ClosePostgresConnection(conn *pgx.Conn) {
	if conn == nil {
		return
	}
	ctx := context.Background()
	conn.Close(ctx)
}
