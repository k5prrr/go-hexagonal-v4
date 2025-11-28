package api

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type AuthRequest struct {
	Phone int64 `json:"phone"`
	Code  int64 `json:"code"`
}
func (r *Router) loginAuthCode(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	request := AuthRequest{}
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"err": "Invalid request"})

		return
	}

	ctx := req.Context()

	phoneStr := strconv.FormatInt(request.Phone, 10)

	if request.Code == 0 {
		err = r.useCase.SendAuthCode(ctx, phoneStr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"err": "Failed SendAuthCode"})

			return
		}
		/*json.NewEncoder(w).Encode(map[string]string{
			"status": "success",
		})*/
		w.Write([]byte(`{}`))
	}

	codeStr := strconv.FormatInt(request.Code, 10)
	id, token, err := r.useCase.CheckAuthCode(ctx, phoneStr, codeStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"err": "Failed CheckAuthCode"})

		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"id": id,
		"token": token,
	})
}


/*
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
	json.NewEncoder(w).Encode(map[string]int{
		"code": 132, // или "message": "SMS sent"
	})
}*/