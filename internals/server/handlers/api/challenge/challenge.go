package challenges

import (
	"encoding/json"
	"fmt"
	"net/http"
	db_challenges "soceng-ai/database/tables/db_challenges"
	challenge_structs "soceng-ai/internals/server/handlers/api/challenge/challenge_structs"
	authentification "soceng-ai/internals/server/handlers/authentification"
)

type Create_challenge_request struct {
	Challenge challenge_structs.Challenge `json:"challenge"`
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
	fmt.Fprintf(w, "soceng-ai/internals/server/handlers/api/challenge/challenges.go:Create\n")
	request := Create_challenge_request{}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err_str := db_challenges.Create_challenge(request.Challenge, r, w)
	if err_str != "" {
		http.Error(w, err_str, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"message":   "Request processed successfully",
		"challenge": request.Challenge.Title,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
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
