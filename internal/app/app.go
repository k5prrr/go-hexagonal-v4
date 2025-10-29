package app

import (
	"context"
	//adapterhttp "app/internal/app/adapter/http"
	"app/pkg/server"
)

type App struct {
	di     *dependencyInjection
	server *server.Server
	config *AppConfig
}
type AppConfig struct {
	HttpAddr   string
	PGUser     string
	PGPassword string
	PGDb       string
	PGPort     string
	PGHost     string
}

func New(config *AppConfig) (*App, error) {
	a := &App{
		di:     NewDependencyInjection(config),
		config: config,
	}
	err := a.initServer()
	if err != nil {

		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	return a.server.Run(ctx)
}

// Инициализация сервера
func (a *App) initServer() error {
	a.server = server.New(a.config.HttpAddr, a.di.Router(), 0)

	return nil
}
