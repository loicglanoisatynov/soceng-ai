package logging

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	database "soceng-ai/database"
	db_cookies "soceng-ai/database/tables/db_cookies"
	"soceng-ai/database/tables/db_profiles"
	db_users "soceng-ai/database/tables/db_users"
	"soceng-ai/internals/utils/prompts"
	"time"
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
	return id != -1
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
		prompts.Prompts_server(time.Now(), prompts.Error+"Error getting user: "+err.Error())
		http.Error(w, "Internal Error", http.StatusInternalServerError)
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

func Cookies_relevant(cookies []*http.Cookie) string {
	// for _, cookie := range cookies {
	// 	if cookie.Name != "socengai-username" && cookie.Name != "socengai-auth" {
	// 		return "Needed cookies : socengai-username & socengai-auth"
	// 	}
	// }
	return "OK"
}

func Is_admin(r *http.Request) bool {
	cookie, err := r.Cookie("socengai-auth")
	if err != nil {
		fmt.Println("Error getting cookie: ", err)
		return false
	}
	username_cookie, err := r.Cookie("socengai-username")
	if err != nil {
		fmt.Println("Error getting username cookie: ", err)
		return false
	}
	if !IsCookieValid(username_cookie.Value, cookie.Value) {
		fmt.Println("Invalid cookie")
		return false
	}
	user_id := db_users.Get_user_id_by_username_or_email(username_cookie.Value)
	if user_id == -1 {
		fmt.Println("User not found")
		return false
	}
	return db_users.Is_admin(user_id)
}

func Get_cookie_value(r *http.Request, cookie_name string) string {
	for _, cookie := range r.Cookies() {
		if cookie.Name == cookie_name {
			return cookie.Value
		}
	}
	return "Error : cookie not found"
}

func Get_user_info(r *http.Request) (string, map[string]interface{}) {
	username_cookie, err := r.Cookie("socengai-username")
	if err != nil {
		return "Error getting username cookie: " + err.Error(), nil
	}
	auth_cookie, err := r.Cookie("socengai-auth")
	if err != nil {
		return "Error getting auth cookie: " + err.Error(), nil
	}

	if !IsCookieValid(username_cookie.Value, auth_cookie.Value) {
		return "Invalid cookie", nil
	}

	user_info, err := db_users.Get_user(database.Get_DB(), "username", username_cookie.Value)
	if err != nil {
		return "Error getting user info: " + err.Error(), nil
	}

	profile, err := db_profiles.Get_profile(user_info.Username)
	if err != nil {
		return "Error getting user biography: profile not found or database error: " + err.Error(), nil
	}

	data := map[string]interface{}{
		"username":  user_info.Username,
		"email":     user_info.Email,
		"admin":     user_info.Is_admin,
		"biography": profile.Biography,
	}

	return "OK", data
}
