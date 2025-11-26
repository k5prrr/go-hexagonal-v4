package database

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

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
	Pool() *pgxpool.Pool
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, query string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	Time() string
	ClearName(text string) string
	ClearText(text string) string
	Close()
}

type ICRUD interface {
	Add(ctx context.Context, entity *any) (int64, error)

	Get(ctx context.Context, id int64) (*any, error)
	GetBy(ctx context.Context, filterKey, filterValue string) (*any, error)
	GetByInt(ctx context.Context, filterKey string, filterValue int64) (*any, error)

	List(ctx context.Context, offset, limit int64) ([]any, error)
	ListBy(ctx context.Context, filterKey, filterValue string, offset, limit int64) ([]any, error)

	Update(ctx context.Context, id int64, entity *any) error
	UpdateBy(ctx context.Context, filterKey, filterValue string, entity *any) error
	UpdateColumn(ctx context.Context, id int64, key, value string) error

	Delete(ctx context.Context, id int64, soft bool) error
	DeleteBy(ctx context.Context, filterKey, filterValue string, soft bool) error
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
func (d *pgxDB) Pool() *pgxpool.Pool {
	return d.pool
}

func (d *pgxDB) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	if d.pool == nil {
		return nil, errors.New("database pool is not initialized")
	}
	return d.pool.Query(ctx, query, args...)
}

func (d *pgxDB) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	return d.pool.Exec(ctx, query, args...)
}

func (d *pgxDB) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	return d.pool.QueryRow(ctx, query, args...)
}

func (d *pgxDB) Begin(ctx context.Context) (pgx.Tx, error) {
	return d.pool.Begin(ctx)
}

func (d *pgxDB) Time() string {
	currentTime := time.Now()

	return currentTime.Format("2006-01-02 15:04:05")
}

func (d *pgxDB) ClearName(text string) string {
	text = strings.TrimSpace(text)
	re := regexp.MustCompile(`[^a-zA-Z0-9_]`)
	text = re.ReplaceAllString(text, "")

	if utf8.RuneCountInString(text) > 255 {
		runes := []rune(text)
		text = string(runes[:255])
	}

	return text
}

func (d *pgxDB) ClearText(text string) string {
	text = strings.TrimSpace(text)

	repl := strings.NewReplacer(
		"\r", "-",
		`"`, "-",
		"'", "-",
		"`", "-",
	)
	text = repl.Replace(text)

	re := regexp.MustCompile(`[\x00\n\r\x1A"'\\]`)

	return re.ReplaceAllStringFunc(text, func(s string) string {
		return `\` + s
	})
}

func (d *pgxDB) Close() {
	if d.pool != nil {
		d.pool.Close()
	}
}
