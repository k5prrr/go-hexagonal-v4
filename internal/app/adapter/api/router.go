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
}

func New(useCase port.IUseCase, path string) *Router {
	r := &Router{
		mux:     http.NewServeMux(),
		useCase: useCase,
		path:    path,
	}
	r.init()
	return r
}

func (r *Router) init() {
	/*	r.mux.Handle(
		fmt.Sprintf("%sdoc/", r.path),
		http.FileServer(http.Dir("./static/")),
	)*/
	docPath := fmt.Sprintf("%sdoc/", r.path)
	// 1. Редирект /api/v2/doc → /api/v2/doc/
	r.mux.HandleFunc(fmt.Sprintf("%sdoc", r.path), func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, docPath, http.StatusMovedPermanently)
	})

	// 2. Файловый сервер с отрезанием префикса
	r.mux.Handle(docPath,
		http.StripPrefix(docPath,
			http.FileServer(http.Dir("./static/")),
		),
	)

	r.testRouter()
	r.authRouter()
	r.telegramRouter()
}

func (r *Router) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) {
	r.mux.HandleFunc(fmt.Sprintf("%s%s", r.path, path), f)
}

// ServeHTTP делает Router совместимым с http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
