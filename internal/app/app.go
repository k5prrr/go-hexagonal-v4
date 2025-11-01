package app

import (
	"app/internal/app/worker/telegram"
	"context"
	//adapterhttp "app/internal/app/adapter/http"
	"app/pkg/server"
)

type App struct {
	di        *dependencyInjection
	server    *server.Server
	config    *AppConfig
	telegramW *telegram.TelegramWorker
}
type AppConfig struct {
	HttpAddr string

	PGUser     string
	PGPassword string
	PGDb       string
	PGPort     string
	PGHost     string

	TGToken      string
	TGTimeoutSec int
}

func New(config *AppConfig) (*App, error) {
	a := &App{
		di:     NewDependencyInjection(config),
		config: config,
	}
	a.server = server.New(a.config.HttpAddr, a.di.Router(), 0)
	a.telegramW = telegram.NewTelegramWorker(a.config.TGToken, a.di.UseCase())

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	a.telegramW.Start(ctx, a.config.TGTimeoutSec)
	return a.server.Run(ctx)
}

// Инициализация сервера
