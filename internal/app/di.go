// Dependency injection / Внедрение зависимостей
package app

import (
	"app/internal/app/adapter/api"
	"app/internal/app/core/port"
	"app/internal/repository/postgres"
	"app/pkg/env"
	"context"
	"database/sql"
	"time"
)

type dependencyInjection struct {
	router   *api.Router
	useCase  *port.IUseCase
	repoAuth *port.IAuthRepo
	env      *env.Env
	db       *pgx.Conn
}

func NewDependencyInjection() *dependencyInjection {
	return &dependencyInjection{
		env: env.New(""),
	}
}

func (d *dependencyInjection) Router() *api.Router {
	if d.router == nil {
		d.router = api.NewRouter(d.UseCase())
	}
	return d.router
}
func (d *dependencyInjection) UseCase() *port.IUseCase {
	if d.useCase == nil {
		//d.useCase = api.NewRouter(d.UseCase)
	}
	return d.useCase
}
func (d *dependencyInjection) RepoAuth() *port.IAuthRepo {
	if d.repoAuth == nil {
		d.repoAuth = postgres.NewRepoAuth(d.DB())
	}
	return d.repoAuth
}
func (d *dependencyInjection) DB() *pgx.Conn {
	if d.db == nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		conn, err := pgx.Connect(ctx, d.env.Get("DATABASE_URL", "postgres://user:pass@localhost/dbname?sslmode=disable"))
		if err != nil {
			panic("failed to connect to database: " + err.Error())
		}
		d.db = conn
	}
	return d.db
}
