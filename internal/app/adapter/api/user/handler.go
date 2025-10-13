//тут и хендлеры

//и роуты

package user

import (
	"app/internal/domain/user"
	"encoding/json"
	"fmt"
	"net/http"
)

type Handlers struct {
	UserService user.AnyUserService
}

func NewHandlers(userService user.AnyUserService) *Handlers {
	return &Handlers{
		UserService: userService,
	}
}
func (h *Handlers) TestSpeed(w http.ResponseWriter, r *http.Request) {
	//http.Error(w, err.Error(), http.StatusInternalServerError)
	//fmt.Fprintf(w, "Index %s\nURL %s", text, r.URL.String())
	fmt.Fprintf(w, "TestSpeed executed successfully")
}
func (h *Handlers) Users(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		// Получение всех пользователей
		users, err := h.UserService.AllUsers()
		if err != nil {
			http.Error(w, JSONError("Failed to fetch users: "+err.Error()), http.StatusInternalServerError)
			return
		}
		_ = json.NewEncoder(w).Encode(users)

	case http.MethodPost:
		// Создание нового пользователя из тела запроса
		var newUser user.User
		if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
			http.Error(w, JSONError("Invalid request body"), http.StatusBadRequest)
			return
		}

		uid, err := h.UserService.CreateUser(&newUser)
		if err != nil {
			http.Error(w, JSONError("Failed to create user: "+err.Error()), http.StatusInternalServerError)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"uid": uid})

	case http.MethodPut:
		// Обновление пользователя по UID из query параметра
		uid := r.URL.Query().Get("uid")
		if uid == "" {
			http.Error(w, JSONError("Missing 'uid' parameter"), http.StatusBadRequest)
			return
		}

		existingUser, err := h.UserService.UserByUid(uid)
		if err != nil {
			http.Error(w, JSONError("User not found: "+err.Error()), http.StatusNotFound)
			return
		}

		// Читаем новые данные из тела запроса и обновляем
		if err := json.NewDecoder(r.Body).Decode(existingUser); err != nil {
			http.Error(w, JSONError("Invalid request body"), http.StatusBadRequest)
			return
		}

		if err := h.UserService.UpdateUser(existingUser); err != nil {
			http.Error(w, JSONError("Failed to update user: "+err.Error()), http.StatusInternalServerError)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "User updated"})

	case http.MethodDelete:
		// Удаление пользователя по UID из query параметра
		uid := r.URL.Query().Get("uid")
		if uid == "" {
			http.Error(w, JSONError("Missing 'uid' parameter"), http.StatusBadRequest)
			return
		}

		if err := h.UserService.DeleteByUid(uid); err != nil {
			http.Error(w, JSONError("Failed to delete user: "+err.Error()), http.StatusInternalServerError)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "User deleted"})

	default:
		http.Error(w, JSONError("Method not allowed"), http.StatusMethodNotAllowed)
	}
}

func JSONError(message string) string {
	// return fmt.Sprintf(`{"error": "%s"}`, message)
	type errResponse struct {
		Error string `json:"error"`
	}
	response, _ := json.Marshal(errResponse{Error: message})
	return string(response)
}

/*
type AnyUserService interface {
	UserByUid(uid string) (*User, error)
	CreateUser(user *User) (string, error)
	UpdateUser(user *User) error
	DeleteByUid(uid string) error
	AllUsers() ([]*User, error)
}
type User struct {
	Uid string

	FamilyName string
	Name       string
	MiddleName string

	BirthDate      time.Time
	Phone          string
	Email          string
	PhoneConfirmed bool
	EmailConfirmed bool

	LastLogin time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time

	Description  string
	PasswordHash string
	KeyApi       string
}
*/
