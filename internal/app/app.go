package app

import (
	"context"
	//adapterhttp "app/internal/app/adapter/http"
	"app/pkg/server"
)

const (
	httpAddr = ":8080"
)

type App struct {
	di     *dependencyInjection
	server *server.Server
}

func New() (*App, error) {
	a := &App{
		di: NewDependencyInjection(),
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
	a.server = server.New(httpAddr, a.di.Router(), 0)

	return nil
}
