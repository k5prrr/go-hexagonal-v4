package api

import (
	"app/internal/app/core/port"
	"encoding/json"
	"net/http"
)

type Router struct {
	mux     *http.ServeMux
	useCase port.IUseCase
	bot     string
}

func New(useCase port.IUseCase, bot string) *Router {
	r := &Router{
		mux:     http.NewServeMux(),
		useCase: useCase,
		bot:     bot,
	}
	r.init()
	return r
}
func (r *Router) init() {
	r.HandleFunc("/api/bot", r.showBot)
	r.HandleFunc("/api/test1", r.test1)

	r.HandleFunc("/api/loginAuthCode", r.loginAuthCode)


	r.mux.Handle("/",
		http.StripPrefix("/",
			http.FileServer(http.Dir("./static/")),
		),
	)
	r.mux.Handle("/api/doc/",
		http.StripPrefix("/api/doc/",
			http.FileServer(http.Dir("./docs/static/")),
		),
	)

}
func (r *Router) showBot(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{
		"bot": r.bot,
	})
}

func (r *Router) test1(w http.ResponseWriter, req *http.Request) {
	user, err := r.useCase.Test1()
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"err": err,
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"user": user,
	})
}

// ---


// Прокидываем для удобства
func (r *Router) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) {
	r.mux.HandleFunc(path, f)
}
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
