package api

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (r *Router) authRouter() {
	r.HandleFunc("registration/linkcheckphone", r.handleLinkCheckPhone)
	r.HandleFunc("addclient/linkcheckphone", r.handleLinkCheckPhone)
	r.HandleFunc("sendpassword", r.sendPassword)
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

type PhoneRequest struct {
	Phone string `json:"phone"`
}

func (r *Router) sendPassword(w http.ResponseWriter, req *http.Request) {
	request := PhoneRequest{}
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	ctx := req.Context()

	phone, err := strconv.ParseInt(request.Phone, 10, 64)
	if err != nil {
		http.Error(w, "Phone not num", http.StatusInternalServerError)
		return
	}

	err = r.useCase.SendPasswordByPhone(ctx, phone)
	if err != nil {
		http.Error(w, "Failed send password", http.StatusInternalServerError)
		return
	}

	// Например, вернуть код (в реальности — отправить SMS, а клиенту — UUID или токен)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{
		"ok": 1, // или "message": "SMS sent"
	})
}
