package api

import (
	"net/http"
	"soceng-ai/database/tables/db_cookies"
	challenge "soceng-ai/internals/server/handlers/api/challenge"
	authentification "soceng-ai/internals/server/handlers/authentification"
)

func Challenge_handler(w http.ResponseWriter, r *http.Request) {

	if process_cookies(w, r) == "error" {
		return
	}

	switch r.Method {
	case "POST":
		challenge.Create(w, r)
	// case "GET":
	// 	challenge.Read(w, r)
	case "PUT":
		challenge.Update(w, r)
	// case "DELETE":
	// 	challenge.Delete(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

func process_cookies(w http.ResponseWriter, r *http.Request) string {
	e := "error"
	cookies := r.Cookies()
	if len(r.Cookies()) < 2 {
		http.Error(w, "Missing cookie.\n", http.StatusUnauthorized)
		return e
	} else if !authentification.Cookies_relevant(cookies, w) {
		return e
	}

	username_cookie, err := r.Cookie("socengai-username")
	if err != nil {
		http.Error(w, "Error getting username cookie: "+err.Error(), http.StatusInternalServerError)
		return e
	}
	auth_cookie, err := r.Cookie("socengai-auth")
	if err != nil {
		http.Error(w, "Error getting auth cookie: "+err.Error(), http.StatusInternalServerError)
		return e
	}

	if username_cookie.Value == "" || auth_cookie.Value == "" {
		http.Error(w, "Cookie empty.\n", http.StatusUnauthorized)
		return e
	}

	if !db_cookies.Is_cookie_valid(username_cookie.Value, auth_cookie.Value) {
		http.Error(w, "Invalid cookie.\n", http.StatusUnauthorized)
		return e
	}
	return ""
}
