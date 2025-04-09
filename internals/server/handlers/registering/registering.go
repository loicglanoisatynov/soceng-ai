package registering

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	database "soceng-ai/database"
	db_users "soceng-ai/database/tables/db_users"
	"soceng-ai/internals/server/handlers/logging"
	"strings"
)

type Registering_request struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Delete_user_request struct {
	Username string `json:"username"`
}

type Registering_response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

func Register_user(w http.ResponseWriter, r *http.Request) {
	request := Registering_request{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Error decoding JSON: %s\n"+err.Error(), http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %s\n", err.Error())
		return
	}
	valid, err := is_register_request_valid(request)

	if !valid {
		http.Error(w, "Invalid request : "+err, http.StatusBadRequest)
	}

	user := db_users.User{
		Username: request.Username,
		Email:    request.Email,
		Password: request.Password,
	}

	if err := db_users.Create_user(database.Get_DB(), user); err != nil {
		http.Error(w, "Failed to create user : "+err.Error(), http.StatusInternalServerError)
		return
	}

	cookie := logging.IssueCookie(user.Username)
	http.SetCookie(w, &http.Cookie{
		Name:  "socengai-username",
		Value: user.Username,
	})
	http.SetCookie(w, &http.Cookie{
		Name:  "socengai-auth",
		Value: cookie,
	})
}

func is_register_request_valid(request Registering_request) (bool, string) {
	if valid, msg := isUserValid(request.Username); !valid {
		return false, msg
	} else if valid, msg := isValidEmail(request.Email); !valid {
		return false, msg
	} else if valid, msg := IsValidPassword(request.Password); !valid {
		return false, msg
	}
	return true, ""
}

func Delete_user(w http.ResponseWriter, r *http.Request) {
	request := Delete_user_request{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Error decoding JSON: %s\n"+err.Error(), http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %s\n", err.Error())
		return
	}

	cookie, err := r.Cookie("socengai-auth")
	if err != nil {
		http.Error(w, "Cookie not found", http.StatusUnauthorized)
		return
	}

	if !logging.IsCookieValid(request.Username, cookie.Value) {
		http.Error(w, "Invalid cookie", http.StatusUnauthorized)
		return
	}

	user_id := db_users.Get_user_id_by_username_or_email(request.Username)
	if user_id == -1 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if !logging.Delete_cookie(request.Username, cookie.Value) {
		http.Error(w, "Failed to delete cookie", http.StatusInternalServerError)
		return
	}

	if err := db_users.Delete_user(database.Get_DB(), user_id); err != nil {
		http.Error(w, "Failed to delete user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "socengai-username",
		Value:  "",
		MaxAge: -1,
	})
	http.SetCookie(w, &http.Cookie{
		Name:   "socengai-auth",
		Value:  "",
		MaxAge: -1,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := Registering_response{
		Status:  true,
		Message: "User deleted successfully",
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("User deleted successfully"))
}

func isUserValid(username string) (bool, string) {
	user, err := db_users.Get_user(database.Get_DB(), "username", username)
	if err == nil && user.ID != 0 {
		return false, "Username already exists"
	} else if len(username) < 1 || len(username) > 30 {
		return false, "Username length must be between 1 and 30 characters"
	} else if strings.ContainsAny(username, "!@#$%^&*()_+[]{}|;':\",.<>?/~`") {
		return false, "Username cannot contain special characters"
	}
	return true, ""
}

func isValidEmail(email string) (bool, string) {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)
	if len(email) < 5 || len(email) > 50 {
		return false, "Email length must be between 5 and 50 characters"
	} else if !re.MatchString(email) {
		return false, "Invalid email format"
	} else if strings.Contains(email, " ") {
		return false, "Email cannot contain spaces"
	} else if user, err := db_users.Get_user(database.Get_DB(), "email", email); err == nil && user.ID != 0 {
		return false, "Email already exists"
	}

	return true, ""
}

func IsValidPassword(password string) (bool, string) {
	if len(password) < 8 || len(password) > 50 {
		return false, "Password length must be between 8 and 50 characters"
	}
	if !strings.ContainsAny(password, "0123456789") {
		return false, "Password must contain at least one digit"
	}
	if !strings.ContainsAny(password, "!@#$%^&*()_+") {
		return false, "Password must contain at least one special character"
	}
	if strings.Contains(password, " ") {
		return false, "Password cannot contain spaces"
	}
	return true, ""
}
