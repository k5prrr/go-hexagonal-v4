package api

import (
	"encoding/json"
	"net/http"
)

func (r *Router) authRouter() {
	r.HandleFunc("registration/linkcheckphone", r.handleLinkCheckPhone)
	r.HandleFunc("addclient/linkcheckphone", r.handleLinkCheckPhone)
	r.HandleFunc("sendAuthCode", r.sendAuthCode)
	r.HandleFunc("bot", r.showBot)
}

func (r *Router) showBot (w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{
		"bot": r.bot,
	})
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
	Phone int64 `json:"phone"`
}

func (r *Router) sendAuthCode(w http.ResponseWriter, req *http.Request) {
	request := PhoneRequest{}
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)

		return
	}

	ctx := req.Context()

	err = r.useCase.SendAuthCode(ctx, request.Phone)
	if err != nil {
		http.Error(w, "Failed SendAuthCode", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	/*json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
	})*/
	w.Write([]byte(`{}`))
}
