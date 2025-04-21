package logging

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	database "soceng-ai/database"
	db_cookies "soceng-ai/database/tables/db_cookies"
	db_users "soceng-ai/database/tables/db_users"
)

type Login_request struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Logout_request struct {
	Username string `json:"username"`
	Cookie   string `json:"cookie"`
}

type Login_response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type Logout_response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

func randString(length int) string {
	bytes := make([]byte, length/2)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func IssueCookie(identifier string) string {
	cookie := randString(32)
	id := db_users.Get_user_id_by_username_or_email(identifier)
	err := db_cookies.Register_cookie(id, cookie)
	if err != nil {
		fmt.Println("Error registering cookie in database: ", err)
		return ""
	}
	return cookie
}

func IsCookieValid(username string, cookie string) bool {
	id := db_cookies.Get_user_id_by_cookie(username, cookie)
	if id == -1 {
		return false
	}
	return true
}

func Delete_cookie(username string, cookie string) bool {
	id := db_cookies.Get_user_id_by_cookie(username, cookie)

	fmt.Println("ID: ", id)
	fmt.Println("Cookie: ", cookie)

	if id == -1 {
		return false
	}
	err := db_cookies.Delete_cookie(id, cookie)
	if err != nil {
		fmt.Println("Error deleting cookie from database: ", err)
		return false
	}
	return true
}

func Login(w http.ResponseWriter, r *http.Request) {
	request := Login_request{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Error decoding JSON: %s\n"+err.Error(), http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %s\n", err.Error())
		return
	}

	user, err := db_users.Get_user(database.Get_DB(), "username", request.Username)
	if err != nil {
		http.Error(w, "Error getting user: %s\n"+err.Error(), http.StatusInternalServerError)
		return
	}

	if user.Password != request.Password {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	cookie := IssueCookie(user.Username)
	if cookie == "" {
		http.Error(w, "Error issuing cookie", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "socengai-auth",
		Value: cookie,
	})
	http.SetCookie(w, &http.Cookie{
		Name:  "socengai-username",
		Value: user.Username,
	})

	response := Login_response{
		Status:  true,
		Message: "Login successful",
	}
	json.NewEncoder(w).Encode(response)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	request := Logout_request{}
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

	if !IsCookieValid(request.Username, cookie.Value) {
		http.Error(w, "Invalid cookie", http.StatusUnauthorized)
		return
	}

	if !Delete_cookie(request.Username, cookie.Value) {
		http.Error(w, "Error deleting cookie", http.StatusInternalServerError)
		return
	}

	response := Logout_response{
		Status:  true,
		Message: "Logout successful",
	}
	json.NewEncoder(w).Encode(response)
}
