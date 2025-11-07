package app

import (
	"context"
)

type App struct {
	di *dependencyInjection
}

func New() *App {
	return &App{
		di: NewDependencyInjection(),
	}
}

func (a *App) Run(ctx context.Context) error {
	conf := a.di.Conf()

	telegramWorker := a.di.TelegramWorker()
	telegramWorker.Run(ctx, conf.Int("TELEGRAM_TIMEOUT", 4))

	server := a.di.Server()
	err := server.Run(ctx) // block
	telegramWorker.Stop()

	return err
}
