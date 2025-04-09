package profiles_handling

import (
	"encoding/json"
	"net/http"
	"soceng-ai/database/tables/db_cookies"
	db_profiles "soceng-ai/database/tables/db_profiles"
	db_users "soceng-ai/database/tables/db_users"
)

type Profile struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Biography string `json:"biography"`
	Avatar    string `json:"avatar"`
}

type Edit_profile_request struct {
	Username  string `json:"username"`
	Biography string `json:"biography"`
	Avatar    string `json:"avatar"`
}

// Modification du profil utilisateur
func Edit_profile(w http.ResponseWriter, r *http.Request) {
	request := Edit_profile_request{}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request.\nData needed : '{\"username\":\"<username>\", \"email\":\"<email>\", \"password\":\"<password>\", \"biography\":\"<biography>\", \"avatar\":\"<avatar>\"}'\n", http.StatusBadRequest)
		return
	}

	username_cookie, err := r.Cookie("socengai-username")
	auth_cookie, err := r.Cookie("socengai-auth")

	if username_cookie.Value == "" || auth_cookie.Value == "" {
		http.Error(w, "Missing cookie.\n", http.StatusUnauthorized)
		return
	}

	if !db_cookies.Is_cookie_valid(username_cookie.Value, auth_cookie.Value) {
		http.Error(w, "Invalid cookie.\n", http.StatusUnauthorized)
		return
	}

	if !db_profiles.Does_profile_exist(request.Username) {
		db_profiles.Create_profile(request.Username)
	}

	user := db_profiles.Db_profile{
		Username:  request.Username,
		Biography: request.Biography,
		Avatar:    request.Avatar,
	}

	if db_profiles.Update_profile(username_cookie.Value, user) != nil {
		http.Error(w, "Failed to update profile.\n", http.StatusInternalServerError)
		return
	}

	// Changer le cookie de l'utilisateur
	if username_cookie.Value != request.Username {
		db_cookies.Delete_cookie(db_users.Get_user_id_by_username_or_email(username_cookie.Value), auth_cookie.Value)
		db_cookies.Register_cookie(db_users.Get_user_id_by_username_or_email(request.Username), auth_cookie.Value)
	}
	w.WriteHeader(http.StatusOK)
}
