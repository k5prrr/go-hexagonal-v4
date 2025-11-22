// Dependency injection / Внедрение зависимостей
package app

import (
	"app/internal/app/adapter/api"
	"app/internal/app/core/port"
	"app/internal/app/core/service"
	"app/internal/app/core/usecase"
	"app/internal/app/worker/telegramW"
	"app/internal/repository/postgres"
	"app/migration"
	"app/pkg/database"
	"app/pkg/env"
	"app/pkg/server"
	"app/pkg/telegram"
	"fmt"
)

type dependencyInjection struct {
	conf           *env.Env
	db             database.IDB
	migration      *migration.Migration
	repo           port.IRepo
	repoUser       port.IRepoUser
	tg             *telegram.Telegram //port.Itg
	service        *service.Service
	useCase        port.IUseCase
	router         *api.Router
	server         *server.Server
	telegramWorker *telegramW.TelegramWorker
}

func NewDependencyInjection() *dependencyInjection {
	return &dependencyInjection{}
}
func (d *dependencyInjection) Conf() *env.Env {
	if d.conf == nil {
		d.conf = env.New(".env")
	}
	return d.conf
}
func (d *dependencyInjection) DB() database.IDB {
	if d.db == nil {
		conf := d.Conf()

		db, err := database.New(&database.DBConfig{
			User:     conf.Get("POSTGRES_USER", ""),
			Password: conf.Get("POSTGRES_PASSWORD", ""),
			Name:     conf.Get("POSTGRES_DB", ""),
			Host:     conf.Get("APP_POSTGRES_HOST", "localhost"),
			Port:     conf.Get("APP_POSTGRES_PORT", "15432"),
		})
		if err != nil {
			panic(fmt.Sprintf("failed to initialize database: %v", err))
		}
		d.db = db
	}
	return d.db
}

func (d *dependencyInjection) Migration() *migration.Migration {
	if d.migration == nil {
		db := d.DB()
		d.migration = migration.New(db.Pool())
	}
	return d.migration
}
func (d *dependencyInjection) Repo() port.IRepo {
	if d.repo == nil {
		d.repo = postgres.New(d.DB())
	}
	return d.repo
}
func (d *dependencyInjection) RepoUser() port.IRepoUser {
	if d.repoUser == nil {
		d.repoUser = postgres.NewRepoUser(d.DB())
	}
	return d.repoUser
}
func (d *dependencyInjection) Tg() *telegram.Telegram {
	if d.tg == nil {
		conf := d.Conf()
		d.tg = telegram.New(&telegram.TelegramConfig{
			Token:   conf.Get("TELEGRAM_TOKEN", ""),
			Webhook: false,
		})
	}
	return d.tg
}
func (d *dependencyInjection) Services() *service.Service {
	if d.service == nil {
		d.service = service.New(d.Repo())
	}
	return d.service
}
func (d *dependencyInjection) UseCase() port.IUseCase {
	if d.useCase == nil {
		d.useCase = usecase.New(d.Services(), d.Repo(), d.Tg())
	}
	return d.useCase
}
func (d *dependencyInjection) Router() *api.Router {
	if d.router == nil {
		conf := d.Conf()
		d.router = api.New(
			d.UseCase(),
			conf.Get("APP_API_PATH", "/"),
			conf.Get("TELEGRAM_BOT", "https://t.me/a"),
		)
	}
	return d.router
}
func (d *dependencyInjection) Server() *server.Server {
	if d.server == nil {
		conf := d.Conf()
		d.server = server.New(
			conf.Get("APP_SERVER_ADDRESS", ":8080"),
			d.Router(),
			0,
		)
	}
	return d.server
}
func (d *dependencyInjection) TelegramWorker() *telegramW.TelegramWorker {
	if d.telegramWorker == nil {
		d.telegramWorker = telegramW.New(d.Tg(), d.UseCase())
	}
	return d.telegramWorker
}
