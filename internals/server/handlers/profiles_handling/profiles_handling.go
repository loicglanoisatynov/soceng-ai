package profiles_handling

import (
	"encoding/json"
	"fmt"
	"net/http"
	"soceng-ai/database/tables/db_cookies"
	db_profiles "soceng-ai/database/tables/db_profiles"
	db_users "soceng-ai/database/tables/db_users"
	authentification "soceng-ai/internals/server/handlers/authentification"
	registering "soceng-ai/internals/server/handlers/registering"
	"soceng-ai/internals/utils/prompts"
	"time"
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

// @Summary		Éditer le profil de l'utilisateur (nom d'utilisateur, biographie, avatar)
// @Description	Éditer le profil de l'utilisateur avec les informations fournies
// @Tags			profiles, users, authentication, edit, profile, avatar, biography, username
// @Accept			json
// @Produce		json
// @Param			username	body		string	true	"Nom d'utilisateur de l'utilisateur"
// @Param			biography	body		string	true	"Biographie de l'utilisateur"
// @Param			avatar		body		string	true	"Avatar de l'utilisateur (URL ou chemin d'accès)"
// @Success		200			{object}	Response	"Profile updated successfully"
// @Failure		400			{string}	string	"Bad Request"
// @Failure		401			{string}	string	"Unauthorized"
// @Failure 405 		{string}	string	"Method Not Allowed"
// @Failure		500			{string}	string	"Internal Server Error"
// @Security		socengai-username
// @Security		socengai-auth
// @Router			/edit-profile [put]
func Edit_profile(w http.ResponseWriter, r *http.Request) {
	request := Edit_profile_request{}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request.\nData needed : '{\"username\":\"<username>\", \"biography\":\"<biography>\", \"avatar\":\"<avatar>\"}'\n", http.StatusBadRequest)
		return
	}

	cookies := r.Cookies()
	if len(r.Cookies()) < 2 {
		http.Error(w, "Missing cookie.\n", http.StatusUnauthorized)
		return
	}
	cookies_status := authentification.Cookies_relevant(cookies)
	if cookies_status != "OK" {
		prompts.Prompts_server(time.Now(), prompts.Error+"soceng-ai/internals/server/handlers/profiles_handling/profiles_handling.go:Edit_profile():Error processing cookies: "+cookies_status)
		http.Error(w, "Needed cookies : socengai-username & socengai-auth\n", http.StatusUnauthorized)
		return
	}

	username_cookie, err := r.Cookie("socengai-username")
	auth_cookie, err := r.Cookie("socengai-auth")

	if username_cookie.Value == "" || auth_cookie.Value == "" {
		prompts.Prompts_server(time.Now(), prompts.Error+"soceng-ai/internals/server/handlers/profiles_handling/profiles_handling.go:Edit_profile():Missing cookie: socengai-username or socengai-auth")
		http.Error(w, "Missing cookie.\n", http.StatusUnauthorized)
		return
	}

	if !db_cookies.Is_cookie_valid(username_cookie.Value, auth_cookie.Value) {
		prompts.Prompts_server(time.Now(), prompts.Error+"soceng-ai/internals/server/handlers/profiles_handling/profiles_handling.go:Edit_profile():Invalid cookie for user "+username_cookie.Value)
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

	newCookie := authentification.IssueCookie(request.Username)
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

// @Summary		Éditer les informations sensibles de l'utilisateur (email, mot de passe)
// @Description	Éditer les informations sensibles de l'utilisateur avec les informations fournies
// @Tags			profiles, users, authentication, edit, profile, email, password
// @Accept			json
// @Produce		json
// @Param			email		body		string	true	"Nouvel email de l'utilisateur"
// @Param			password	body		string	true	"Mot de passe actuel de l'utilisateur"
// @Param			newpassword	body	string	false	"Nouveau mot de passe de l'utilisateur"
// @Success		200			{object}	Response	"User updated successfully"
// @Failure		400			{string}	string	"Bad Request"
// @Failure		401			{string}	string	"Unauthorized"
// @Failure		405 		{string}	string	"Method Not Allowed"
// @Failure		500			{string}	string	"Internal Server Error"
// @Security		socengai-username
// @Security		socengai-auth
// @Router			/edit-user [put]
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

	valid, error_status := registering.Is_password_users(username_cookie.Value, request.Password)
	fmt.Println(prompts.Debug + "soceng-ai/internals/server/handlers/profiles_handling/profiles_handling.go:Edit_user():Password validity : " + fmt.Sprint(valid))
	fmt.Println(prompts.Debug + "soceng-ai/internals/server/handlers/profiles_handling/profiles_handling.go:Edit_user():Password received : " + request.Password)

	if !valid {

		http.Error(w, "Invalid password : "+error_status+"\n", http.StatusBadRequest)
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
