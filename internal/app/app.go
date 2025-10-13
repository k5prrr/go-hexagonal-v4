// Тут константы, таймауты
package app

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
	//adapterhttp "app/internal/app/adapter/http"
)

const (
	httpAddr          = ":8080"
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

type App struct {
	diContainer *diContainer
	httpServer  *http.Server
}

func New() (*App, error) {
	a := &App{}
	err := a.initDeps()
	if err != nil {
		return nil, err
	}

	return a, nil
}

// Инициализация зависимостей по очереди
func (a *App) initDeps() error {
	inits := []func() error{
		a.initDI,
		a.initHTTPServer,
	}

	for _, f := range inits {
		err := f()
		if err != nil {
			return err
		}
	}

	return nil
}

// Подключение внедрения зависимостей
func (a *App) initDI() error {
	a.diContainer = NewDiContainer()
	return nil
}

// Инициализация сервера
func (a *App) initHTTPServer() error {
	router := http.NewServeMux()
	router.Handle("/", http.FileServer(http.Dir("./static/")))

	// mini test
	var counter uint64
	// router.HandleFunc("/api/increment", incrementHandler))
	router.HandleFunc("/api/v2/increment", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		v := atomic.AddUint64(&counter, 1)
		resp := map[string]uint64{"value": v}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	a.httpServer = &http.Server{
		Addr:              httpAddr,
		ReadHeaderTimeout: readHeaderTimeout,
		Handler:           router,
	}

	return nil
}

// Запуск сервера
func (a *App) runHTTPServer() error {
	log.Printf("starting http server on %s\n", httpAddr)

	err := a.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
func (a *App) Run(ctx context.Context) error {
	go func() {
		if err := a.runHTTPServer(); err != nil {
			log.Printf("http server listen failed: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Printf("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(ctx, shutdownTimeout)
	defer cancel()

	if err := a.httpServer.Shutdown(ctx); err != nil {
		log.Printf("http server shutdown failed: %v\n", err)
		return err
	}

	log.Printf("http server stopped")

	return nil
}
