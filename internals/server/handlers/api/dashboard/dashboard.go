package dashboard

import (
	"fmt"
	"net/http"
	database "soceng-ai/database"
	db_achievements "soceng-ai/database/tables/db_achievements"
	db_challenges "soceng-ai/database/tables/db_challenges"
	db_challenges_structs "soceng-ai/database/tables/db_challenges/db_challenges_structs"
	db_sessions "soceng-ai/database/tables/db_sessions"
	"soceng-ai/internals/server/handlers/api/dashboard/dashboard_structs"
	authentification "soceng-ai/internals/server/handlers/authentification"
	"soceng-ai/internals/utils/prompts"
)

// Récupère les données de dashboard (donc l'énumération des challenges disponibles)
func Get_dashboard(request *http.Request) (string, dashboard_structs.Dashboard) {
	// Récupérer les données de dashboard
	dashboard_data := dashboard_structs.Dashboard{
		Score:      db_achievements.Get_user_score(authentification.Get_cookie_value(request, "socengai-username")),
		Challenges: db_challenges.Get_available_challenges(authentification.Get_cookie_value(request, "socengai-username")),
	}

	return "OK", dashboard_data
}

func Get_available_challenges(username string) []dashboard_structs.Challenge {
	var challenges []db_challenges_structs.Challenge
	query := "SELECT * FROM challenges WHERE validated = TRUE"

	db := database.Get_DB()
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error getting available challenges:", err)
		return nil
	}
	defer rows.Close()

	// Parcours des résultats de la requête
	// On crée un tableau de challenges
	for rows.Next() {
		var challenge db_challenges_structs.Challenge
		err := rows.Scan(&challenge.ID, &challenge.Title, &challenge.Lore_for_player, &challenge.Lore_for_ai, &challenge.Difficulty, &challenge.Illustration, &challenge.Created_at, &challenge.Updated_at, &challenge.Validated, &challenge.Osint_data)
		if err != nil {
			fmt.Println("Error scanning challenge:", err)
			continue
		}
		challenges = append(challenges, challenge)
	}
	// On crée un tableau de challenges formaté pour le dashboard
	dashboard_challenges := []dashboard_structs.Challenge{}
	for i := 0; i < len(challenges); i++ {
		var dashboard_challenge dashboard_structs.Challenge
		dashboard_challenge.ID = challenges[i].ID
		dashboard_challenge.Name = challenges[i].Title
		dashboard_challenge.Description = challenges[i].Lore_for_player
		dashboard_challenge.Illustration_filename = challenges[i].Illustration
		dashboard_challenge.Status = "available"
		dashboard_challenges = append(dashboard_challenges, dashboard_challenge)
	}

	// Récupère les sessions de jeu de l'utilisateur. Pour toutes les sessions en cours, on change le status du challenge
	// en "in progress" et on le renvoie
	error_status, sessions := db_sessions.Get_sessions_by_username(username)
	if error_status != "OK" {
		fmt.Println(prompts.Error + "soceng-ai/internals/server/handlers/api/dashboard/dashboard.go:Get_available_challenges():Error getting sessions by username: " + error_status)
	}
	for i := 0; i < len(sessions); i++ {
		for j := 0; j < len(dashboard_challenges); j++ {

			if sessions[i].ChallengeID == dashboard_challenges[j].ID {
				dashboard_challenges[j].Status = "in progress"
				break
			}
		}
	}

	// Récupère les sessions de jeu finies (chaque entrée de l'user dans la table achievement). Chaque challenge fini voit son status modifié pour "done"

	return dashboard_challenges
}
