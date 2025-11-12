package app

import (
	"context"
	"fmt"
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

	migration := a.di.Migration()
	if err := migration.Run(ctx); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	telegramWorker := a.di.TelegramWorker()
	go telegramWorker.Run(ctx, conf.Int("TELEGRAM_TIMEOUT", 4))

	server := a.di.Server()
	err := server.Run(ctx) // block
	telegramWorker.Stop()

	return err
}
