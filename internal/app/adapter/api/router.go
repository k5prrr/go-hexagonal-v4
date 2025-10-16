package api

import (
	"app/internal/app/core/port"
	"fmt"
	"net/http"
)

type Router struct {
	mux     *http.ServeMux
	useCase *port.IUseCase
	path    string
}

func NewRouter(useCase *port.IUseCase) *Router {
	r := &Router{
		mux:     http.NewServeMux(),
		useCase: useCase,
		path:    "/api/v2/",
	}
	r.initRouter()
	return r
}

func (r *Router) initRouter() {
	r.mux.Handle(fmt.Sprintf("%sdoc", r.path), http.FileServer(http.Dir("./static/")))
	r.testRouter()
	r.authRouter()
}

func (r *Router) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) {
	r.mux.HandleFunc(fmt.Sprintf("%sincrement", r.path), f)
}

// ServeHTTP делает Router совместимым с http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
