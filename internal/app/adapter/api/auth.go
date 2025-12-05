package api

import (
	"app/internal/app/core/domain"
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
	token, err := r.useCase.CheckAuthCode(ctx, phoneStr, codeStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"err": "Failed CheckAuthCode"})

		return
	}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // true только в HTTPS-окружении (prod)
		SameSite: http.SameSiteLaxMode,
		// MaxAge:   3600 * 24 * 7, // ← добавь, если нужен срок годности
	}
	http.SetCookie(w, cookie)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": token,
	})
}

func (r *Router) currentUserH(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	currentUser, err := r.useCase.CurrentUser(ctx, token)
	json.NewEncoder(w).Encode(map[string]string{
		"bot": r.bot,
	})
}

func (r *Router) currentUser(w http.ResponseWriter, req *http.Request) *domain.UserFull {
	ctx := req.Context()

	cookieToken, err := req.Cookie("token")
	if err != nil {
		return nil
	}

	currentUser, err := r.useCase.CurrentUser(ctx, cookieToken.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"err": "Failed CurrentUser"})
	}

	return currentUser
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
