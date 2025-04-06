package registering

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	database "soceng-ai/database"
	users "soceng-ai/database/tables/users"
	"soceng-ai/internals/server/handlers/logging"
	"strings"
)

type Registering_request struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
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

	user := users.User{
		Username: request.Username,
		Email:    request.Email,
		Password: request.Password,
	}

	if err := users.Create_user(database.Get_DB(), user); err != nil {
		http.Error(w, "Failed to create user : "+err.Error(), http.StatusInternalServerError)
		return
	}

	cookie := logging.IssueCookie(user.Username)
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

func Read_user(w http.ResponseWriter, r *http.Request) {
	request := Registering_request{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := users.Get_user(database.Get_DB(), "username", request.Username)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		fmt.Fprintf(w, "User not found: %s", request.Username)
		return
	}

	response := Registering_response{
		Status:  true,
		Message: "User found",
	}
	json.NewEncoder(w).Encode(response)
	fmt.Fprintf(w, "User found: %s", request.Username)
}

func Update_user(w http.ResponseWriter, r *http.Request) {
	request := Registering_request{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := users.Get_user(database.Get_DB(), "username", request.Username)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if valid, msg := isValidEmail(request.Email); !valid {
		http.Error(w, msg, http.StatusBadRequest)
		return
	} else if valid, msg := IsValidPassword(request.Password); !valid {
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	user.Email = request.Email
	user.Password = request.Password

	if err := users.Update_user(database.Get_DB(), user); err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	response := Registering_response{
		Status:  true,
		Message: "User updated successfully",
	}
	json.NewEncoder(w).Encode(response)
}

func Delete_user(w http.ResponseWriter, r *http.Request) {
	request := Registering_request{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := users.Get_user(database.Get_DB(), "username", request.Username)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if err := users.Delete_user(database.Get_DB(), user.ID); err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	response := Registering_response{
		Status:  true,
		Message: "User deleted successfully",
	}
	json.NewEncoder(w).Encode(response)
}

func isUserValid(username string) (bool, string) {
	user, err := users.Get_user(database.Get_DB(), "username", username)
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
	} else if user, err := users.Get_user(database.Get_DB(), "email", email); err == nil && user.ID != 0 {
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
