package api

import (
	"app/internal/app/core/port"
	"fmt"
	"net/http"
)

type Router struct {
	mux     *http.ServeMux
	useCase port.IUseCase
	path    string
	bot string
}

func New(useCase port.IUseCase, path string, bot string) *Router {
	r := &Router{
		mux:     http.NewServeMux(),
		useCase: useCase,
		path:    path,
		bot:     bot,
	}
	r.init()
	return r
}

func (r *Router) init() {
	r.static()
	r.testRouter()
	r.authRouter()
	r.telegramRouter()

	
}

// Для быстрого добавления по нужному пути
func (r *Router) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) {
	r.mux.HandleFunc(fmt.Sprintf("%s%s", r.path, path), f)
}

// ServeHTTP делает Router совместимым с http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
