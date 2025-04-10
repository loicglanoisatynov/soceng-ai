package handlers

import (
	"encoding/json"
	"net/http"

	"soceng-ai/database"
	"soceng-ai/internals/server/model"
	"soceng-ai/internals/utils/debug"
)

func CreateChallenge(w http.ResponseWriter, r *http.Request) {
	var challenge model.Challenge

	// Décodage du body
	if err := json.NewDecoder(r.Body).Decode(&challenge); err != nil {
		debug.LogError("Erreur de décodage du challenge", err)
		http.Error(w, "Données invalides", http.StatusBadRequest)
		return
	}

	db := database.Get_DB()

	query := `
		INSERT INTO challenges (title, description, flag, difficulty)
		VALUES (?, ?, ?, ?)
	`

	_, err := db.Exec(query, challenge.Title, challenge.Description, challenge.Flag, challenge.Difficulty)
	if err != nil {
		debug.LogError("Erreur lors de l'insertion du challenge", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Challenge créé avec succès",
	})
}
