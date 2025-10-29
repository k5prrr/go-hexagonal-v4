package api

import "net/http"

//import "app/pkg/telegram"

func (r *Router) telegramRouter() {
	r.HandleFunc("telegram", r.telegramInput)
}

func (r *Router) telegramInput(w http.ResponseWriter, req *http.Request) {
	//telegram.New()
}
