package challenges

import (
	"encoding/json"
	"fmt"
	"net/http"
	db_challenges "soceng-ai/database/tables/db_challenges"
)

type Create_challenge_request struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	Illustration string `json:"illustration"`
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
}
