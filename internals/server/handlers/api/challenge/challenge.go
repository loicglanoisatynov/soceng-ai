package challenges

import (
	"encoding/json"
	"fmt"
	"net/http"
	db_challenges "soceng-ai/database/tables/db_challenges"
	authentification "soceng-ai/internals/server/handlers/authentification"
)

type Create_challenge_request struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	Illustration string `json:"illustration"`
}

type Update_challenge_request struct {
	Operation    string `json:"operation"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Illustration string `json:"illustration"`
}

type Validate_challenge_request struct {
	Title string `json:"title"`
}

// Création de challenge. Contrôle l'identité de l'utilisateur (Is_cookie_valid) et la validité des données (Is_data_valid).
// Ensuite : créer le système de gestion des interlocuteurs du challenge.
// Vérifier que le challenge est complet
func Create(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "soceng-ai/internals/server/handlers/api/challenge/challenges.go:Create(w, r)\n")
	request := Create_challenge_request{}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request.\nData needed : '{\"title\":\"<title>\", \"description\":\"<description>\", \"illustration\":\"<illustration>\"}'\n", http.StatusBadRequest)
		return
	}

	err_msg, valid := db_challenges.Is_data_valid(request.Title, request.Description, request.Illustration)
	if !valid {
		http.Error(w, err_msg, http.StatusBadRequest)
		return
	}

	db_challenges.Create_challenge(request.Title, request.Description, request.Illustration)

	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"message":     "Challenge created successfully\nTODO : add entities to the challenge",
		"title":       request.Title,
		"description": request.Description,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("TODO ADD TEXT\n"))
}

func Update(w http.ResponseWriter, r *http.Request) {
	request := Update_challenge_request{}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request.\nData needed : '{\"title\":\"<title>\"}'\n", http.StatusBadRequest)
		return
	}

	if request.Operation == "validate" {
		if !authentification.Is_admin(r) {
			http.Error(w, "Unauthorized\n", http.StatusUnauthorized)
			return
		} else {
			db_challenges.Validate_challenge(request.Title)
			return
		}
	}
}
