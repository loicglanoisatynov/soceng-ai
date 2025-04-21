package profiles_handling

import (
	"encoding/json"
	"net/http"
	"soceng-ai/database/tables/db_cookies"
	db_profiles "soceng-ai/database/tables/db_profiles"
	db_users "soceng-ai/database/tables/db_users"
	logging "soceng-ai/internals/server/handlers/logging"
	registering "soceng-ai/internals/server/handlers/registering"
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

type Response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type Edit_user_request struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	Newpassword string `json:"newpassword"`
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

	newCookie := logging.IssueCookie(request.Username)
	// Changer le cookie de l'utilisateur
	if username_cookie.Value != request.Username {
		db_cookies.Delete_cookie(db_users.Get_user_id_by_username_or_email(username_cookie.Value), auth_cookie.Value)
		db_cookies.Register_cookie(db_users.Get_user_id_by_username_or_email(request.Username), newCookie)
	}

	// Mettre à jour le cookie de l'utilisateur
	http.SetCookie(w, &http.Cookie{
		Name:  "socengai-username",
		Value: request.Username,
	})
	http.SetCookie(w, &http.Cookie{
		Name:  "socengai-auth",
		Value: newCookie,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := Response{
		Status:  true,
		Message: "Profile updated successfully",
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Profile updated successfully\n"))
}

// Même chose, mais pour l'édition des informations sensibles (email, mot de passe), conditionné par le mot de passe
// de l'utilisateur
func Edit_user(w http.ResponseWriter, r *http.Request) {
	request := Edit_user_request{}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request.\nData needed : '{\"email\":\"<email>\", \"password\":\"<password>\", \"newpassword\":\"<newpassword>\"}'\n", http.StatusBadRequest)
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

	valid, _ := registering.IsValidPassword(request.Password)
	if !valid {
		http.Error(w, "Invalid password.\n", http.StatusBadRequest)
		return
	}

	valid, _ = registering.IsValidEmail(request.Email)
	if !valid {
		http.Error(w, "Invalid email.\n", http.StatusBadRequest)
		return
	}

	user_id := db_users.Get_user_id_by_username_or_email(username_cookie.Value)

	if user_id == -1 {
		http.Error(w, "User not found.\n", http.StatusNotFound)
		return
	}

	if db_users.Is_password_valid(user_id, request.Password) == false {
		http.Error(w, "Invalid password.\n", http.StatusUnauthorized)
		return
	}

	if request.Email != "" {
		db_users.Update_email(user_id, request.Email)
	}
	if request.Newpassword != "" {
		db_users.Update_password(user_id, request.Newpassword)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := Response{
		Status:  true,
		Message: "User updated successfully",
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("User updated successfully\n"))
}
