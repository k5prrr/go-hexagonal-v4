package api

import (
	"encoding/json"
	"net/http"
)

func (r *Router) authRouter() {
	r.HandleFunc("registration/linkcheckphone", r.handleLinkCheckPhone)
	r.HandleFunc("addclient/linkcheckphone", r.handleLinkCheckPhone)

	r.HandleFunc("sendAuthCode", r.sendAuthCode)
	r.HandleFunc("loginCode", r.loginCode)
	r.HandleFunc("bot", r.showBot)
}


func (r *Router) handleLinkCheckPhone(w http.ResponseWriter, req *http.Request) {
	// Тут проверка полей и отправка в чистые useCase
	//w.WriteHeader(http.StatusOK)
	//w.Write([]byte("OK"))

	/*	ctx := req.Context()

		code, err := r.useCase.CreateCodeCheckPhone(ctx, "registration")
		if err != nil {
			http.Error(w, "Failed to generate code", http.StatusInternalServerError)

			return
		}
	*/

	// Например, вернуть код (в реальности — отправить SMS, а клиенту — UUID или токен)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{
		"code": 132, // или "message": "SMS sent"
	})
}

type AuthRequest struct {
	Phone int64 `json:"phone"`
	Code  int64 `json:"code"`
}

func (r *Router) sendAuthCode(w http.ResponseWriter, req *http.Request) {
	request := AuthRequest{}
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)

		return
	}

	//ctx := req.Context()

	/*err = r.useCase.SendAuthCode(ctx, request.Phone)
	if err != nil {
		http.Error(w, "Failed SendAuthCode", http.StatusInternalServerError)

		return
	}*/

	w.Header().Set("Content-Type", "application/json")
	/*json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
	})*/
	w.Write([]byte(`{}`))
}

func (r *Router) loginCode(w http.ResponseWriter, req *http.Request) {
	request := AuthRequest{}
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)

		return
	}

	//ctx := req.Context()

	/*id, token, err := r.useCase.LoginCode(ctx, request.Phone, request.Code)
	if err != nil {
		http.Error(w, "Failed LoginCode", http.StatusInternalServerError)

		return
	}*/

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":    0,
		"token": 0,
	})
	//w.Write([]byte(`{}`))
}
