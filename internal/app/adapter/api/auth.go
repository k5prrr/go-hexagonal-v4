package api

import (
	"encoding/json"
	"net/http"
)

func (r *Router) authRouter() {
	r.HandleFunc("registration/linkcheckphone", r.handleLinkCheckPhone)
	r.HandleFunc("addclient/linkcheckphone", r.handleLinkCheckPhone)
}

func (r *Router) handleLinkCheckPhone(w http.ResponseWriter, req *http.Request) {
	// Тут проверка полей и отправка в чистые useCase
	//w.WriteHeader(http.StatusOK)
	//w.Write([]byte("OK"))

	ctx := req.Context()

	code, err := r.useCase.CreateCodeCheckPhone(ctx, "registration")
	if err != nil {
		http.Error(w, "Failed to generate code", http.StatusInternalServerError)
		return
	}

	// Например, вернуть код (в реальности — отправить SMS, а клиенту — UUID или токен)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"code": code, // или "message": "SMS sent"
	})
}
