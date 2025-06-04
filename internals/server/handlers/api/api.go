package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"soceng-ai/database/tables/db_cookies"
	"soceng-ai/database/tables/db_sessions"
	challenge "soceng-ai/internals/server/handlers/api/challenge"
	dashboard "soceng-ai/internals/server/handlers/api/dashboard"
	sessions "soceng-ai/internals/server/handlers/api/sessions"
	authentification "soceng-ai/internals/server/handlers/authentification"
	"soceng-ai/internals/utils/prompts"
	"time"
)

var re = regexp.MustCompile(`^[a-zA-Z0-9]{6}$`)

func Challenge_handler(w http.ResponseWriter, r *http.Request) {

	cookies_status := process_cookies(r)
	if cookies_status != "OK" {
		http.Error(w, "Error processing cookies", http.StatusUnauthorized)
		return
	}

	switch r.Method {

	// Créer le challenge
	case "POST":
		challenge.Create(w, r)

	// Récupérer le challenge
	// case "GET":
	// 	challenge.Read(w, r)

	// Valider le challenge
	case "PUT":
		challenge.Update(w, r)

	// Supprimer le challenge
	// case "DELETE":
	// 	challenge.Delete(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

func Sessions_handler(w http.ResponseWriter, r *http.Request) {
	var error_status string
	var response = http.Response{
		StatusCode: http.StatusBadRequest,
		Status:     "Bad Request",
	}

	cookies_status := process_cookies(r)
	if cookies_status != "OK" {
		http.Error(w, "Error processing cookies: "+cookies_status, http.StatusUnauthorized)
		return
	}

	if r.URL.Path == "/api/sessions/start-challenge" {
		switch r.Method {
		case "POST":
			response = sessions.Start_challenge(r)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	} else {
		if re.MatchString(r.URL.Path[len("/api/sessions/"):]) {
			session_id := r.URL.Path[len("/api/sessions/"):]
			// On vérifie que la session existe
			error_status, response.StatusCode = db_sessions.Check_session_key(r, session_id)
			if error_status != "OK" {
				fmt.Println(prompts.Error + "soceng-ai/internals/server/handlers/api/sessions/sessions.go:Sessions_handler():Error checking session ID: " + error_status)
				http.Error(w, "Error checking session ID: "+error_status, http.StatusBadRequest)
				return
			}
			switch r.Method {
			case "GET":
				sessions.Get_session_data(r, w, session_id)

			case "POST":
				response = sessions.Post_session_data(r, session_id)
			default:

				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		} else {
			http.Error(w, "Invalid session ID", http.StatusBadRequest)
			return
		}
	}

	response.Write(w)
}

func process_cookies(r *http.Request) string {
	cookies := r.Cookies()
	if len(cookies) < 2 {
		return "Missing cookie"
	}
	cookies_status := authentification.Cookies_relevant(cookies)
	if cookies_status != "OK" {
		return cookies_status
	}

	username_cookie, err := r.Cookie("socengai-username")
	if err != nil {
		return "Error getting username cookie : " + err.Error()
	}
	auth_cookie, err := r.Cookie("socengai-auth")
	if err != nil {
		return "Error getting auth cookie : " + err.Error()
	}

	if username_cookie.Value == "" || auth_cookie.Value == "" {
		return "Cookie empty"
	}

	cookies_valid, err_msg := db_cookies.Is_cookie_valid(username_cookie.Value, auth_cookie.Value)
	if !cookies_valid {
		return "Invalid cookie : " + err_msg
	}
	return "OK"
}

// @Summary Dashboard handler
// @Description Gère les requêtes pour récupérer les données du tableau des challenges de la page principale.
// @Tags Dashboard, Challenges
// @Accept json
// @Produce json
// @Success 200 {object} dashboard_structs.Dashboard
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/dashboard [get]
func Dashboard_handler(w http.ResponseWriter, r *http.Request) {

	cookies_status := process_cookies(r)
	if cookies_status != "OK" {
		http.Error(w, "Error processing cookies: "+cookies_status, http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case "GET":
		error_status, data := dashboard.Get_dashboard(r)
		if error_status != "OK" {
			fmt.Println(prompts.Error + "soceng-ai/internals/server/handlers/api/dashboard/dashboard.go:Dashboard_handler():Error getting dashboard data: " + error_status)
			http.Error(w, "Error getting dashboard data: "+error_status, http.StatusInternalServerError)
			return
		}
		responseData, err := json.Marshal(data)
		if err != nil {
			fmt.Println(prompts.Error + "soceng-ai/internals/server/handlers/api/dashboard/dashboard.go:Dashboard_handler():Error marshalling dashboard data: " + err.Error())
			http.Error(w, "Error marshalling dashboard data: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(responseData)
		if err != nil {
			fmt.Println(prompts.Error + "soceng-ai/internals/server/handlers/api/dashboard/dashboard.go:Dashboard_handler():Error writing response: " + err.Error())
			http.Error(w, "Error writing response: "+err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func HelloWorld_handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Hello, World !")
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func NotHelloWorld_handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Not Hello, World !")
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func Check_cookies_handler(w http.ResponseWriter, r *http.Request) {
	cookies := r.Cookies()
	if len(cookies) < 2 {
		prompts.Prompts_server(time.Now(), "soceng-ai/internals/server/handlers/api/check_cookies.go:Check_cookies_handler():Missing cookie")
		http.Error(w, "Missing cookie", http.StatusUnauthorized)
		return
	}

	cookies_status := authentification.Cookies_relevant(cookies)
	if cookies_status != "OK" {
		prompts.Prompts_server(time.Now(), "soceng-ai/internals/server/handlers/api/check_cookies.go:Check_cookies_handler():Error processing cookies: "+cookies_status)
		http.Error(w, "Error processing cookies: "+cookies_status, http.StatusUnauthorized)
		return
	}

	username_cookie, err := r.Cookie("socengai-username")
	if err != nil {
		prompts.Prompts_server(time.Now(), "soceng-ai/internals/server/handlers/api/check_cookies.go:Check_cookies_handler():Error getting username cookie: "+err.Error())
		http.Error(w, "Error getting username cookie: "+err.Error(), http.StatusUnauthorized)
		return
	}
	auth_cookie, err := r.Cookie("socengai-auth")
	if err != nil {
		prompts.Prompts_server(time.Now(), "soceng-ai/internals/server/handlers/api/check_cookies.go:Check_cookies_handler():Error getting auth cookie: "+err.Error())
		http.Error(w, "Error getting auth cookie: "+err.Error(), http.StatusUnauthorized)
		return
	}

	if username_cookie.Value == "" || auth_cookie.Value == "" {
		prompts.Prompts_server(time.Now(), "soceng-ai/internals/server/handlers/api/check_cookies.go:Check_cookies_handler():Cookie empty")
		http.Error(w, "Cookie empty", http.StatusUnauthorized)
		return
	}

	cookie_valid, err_msg := db_cookies.Is_cookie_valid(username_cookie.Value, auth_cookie.Value)
	if !cookie_valid {
		http.Error(w, "Invalid cookie : "+err_msg, http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"status": "OK"}`)
}

func User_info_handler(w http.ResponseWriter, r *http.Request) {
	cookies_status := process_cookies(r)
	if cookies_status != "OK" {
		http.Error(w, "Error processing cookies: "+cookies_status, http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case "GET":
		error_status, data := authentification.Get_user_info(r)
		if error_status != "OK" {
			prompts.Prompts_server(time.Now(), prompts.Error+"soceng-ai/internals/server/handlers/api/user_info.go:User_info_handler() : Error getting user info: "+error_status)
			http.Error(w, "Error getting user info: "+error_status, http.StatusInternalServerError)
			return
		}
		responseData, err := json.Marshal(data)
		if err != nil {
			fmt.Println(prompts.Error + "soceng-ai/internals/server/handlers/api/user_info.go:User_info_handler():Error marshalling user info: " + err.Error())
			http.Error(w, "Error marshalling user info: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(responseData)
		if err != nil {
			fmt.Println(prompts.Error + "soceng-ai/internals/server/handlers/api/user_info.go:User_info_handler():Error writing response: " + err.Error())
			http.Error(w, "Error writing response: "+err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
