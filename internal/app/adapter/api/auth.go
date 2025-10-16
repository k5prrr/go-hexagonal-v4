package api

import "net/http"

func (r *Router) authRouter() {
	r.mux.HandleFunc("/api/v2/registration/linkcheckphone", r.handleLinkCheckPhone)
	r.mux.HandleFunc("/api/v2/addclient/linkcheckphone", r.handleLinkCheckPhone)
}

func (r *Router) handleLinkCheckPhone(w http.ResponseWriter, req *http.Request) {
	// Тут проверка полей и отправка в чистые useCase
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
