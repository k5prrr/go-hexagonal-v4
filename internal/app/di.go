// Dependency injection / Внедрение зависимостей
package app

import (
	"app/internal/app/adapter/api"
	"app/internal/app/core/port"
	"app/internal/app/core/usecase"
	"app/internal/repository/postgres"
	"app/pkg/database"
	"fmt"
)

type dependencyInjection struct {
	router  *api.Router
	useCase port.IUseCase
	config  *AppConfig
	db      database.IDB
	repo    port.IRepo
}

func NewDependencyInjection(config *AppConfig) *dependencyInjection {
	return &dependencyInjection{
		config: config,
	}
}

func (d *dependencyInjection) Router() *api.Router {
	if d.router == nil {
		d.router = api.NewRouter(d.UseCase())
	}
	return d.router
}
func (d *dependencyInjection) UseCase() port.IUseCase {
	if d.useCase == nil {
		d.useCase = usecase.NewUseCase(d.Repo())
	}
	return d.useCase
}
func (d *dependencyInjection) Repo() port.IRepo {
	if d.repo == nil {
		d.repo = postgres.NewRepo(d.DB())
	}
	return d.repo
}

func (d *dependencyInjection) DB() database.IDB {
	if d.db == nil {
		dbConf := &database.DBConfig{
			Name:     d.config.PGDb,
			User:     d.config.PGUser,
			Password: d.config.PGPassword,
			Host:     d.config.PGHost,
			Port:     d.config.PGPort,
		}

		db, err := database.NewDB(dbConf)
		if err != nil {
			panic(fmt.Sprintf("failed to initialize database: %v", err))
		}
		d.db = db
	}
	return d.db
}
