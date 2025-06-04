package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"soceng-ai/database/tables/db_challenges"
	"soceng-ai/database/tables/db_sessions"
	"soceng-ai/internals/server/handlers/api/sessions/sessions_structs"
	"soceng-ai/internals/utils/prompts"
	"strconv"
	"time"
)

var (
	apiKey string = "AIzaSyDVLqTrPlzSI3KvEJHEN58uFQeRPY3rPTU"
)

// Personnage pour le prompt
type PersoData struct {
	Nom                string
	Titre              string
	Organisation       string
	Psychologie        string
	Suspicion          int
	Osint              string
	Document           string
	Contact            string
	MessagePrecedent   string
	MessageUtilisateur string
}

// Structure pour l'appel API Gemini
type Content struct {
	Parts []map[string]string `json:"parts"`
	Role  string              `json:"role"`
}

type RequestBody struct {
	Contents []Content `json:"contents"`
}

// Structure pour parser la réponse JSON de Gemini
type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

// Fonction pour générer le prompt à partir des données du personnage
func GeneratePrompt(data PersoData) string {
	prompt := `Ce message vient d'une API de CTF d'ingénierie sociale, tu dois répondre en fonction des données de jeu passées. Tu as le droit de modifier les données identifiées Altérables, tu as le droit de toucher aux autres. Réponds conformément à ton personnage dans la zone de texte Réponse. Si tu renvoie un niveau de suspicion de 10, tu dois signaler à ton interlocuteur que tu as détecté sa manipulation et que tu alerte la sécurité ou la direction, et qu'il doit quitter les lieux immédiatement ou couper la communication.
"ton nom":"` + data.Nom + `",
"ton titre":"` + data.Titre + `",
"ton organisation":"` + data.Organisation + `",
"ta psychologie":"` + data.Psychologie + `",
"ton niveau de suspicion [ALTERABLE entre 1 et 10]":"` + fmt.Sprint(data.Suspicion) + `",
"tes traces osint":"` + data.Osint + `",`

	if data.Document != "" {
		prompt += "\n\"possede_document\":\"" + data.Document + `",`
	}
	if data.Contact != "" {
		prompt += "\n\"possede_contact\":\"" + data.Contact + `",`
	}

	prompt += `
"ton precedent message":"` + data.MessagePrecedent + `",
"message utilisateur":"` + data.MessageUtilisateur + `"

Renvoie le bloc suivant sans le moindre ajout :
{
"Suspicion (entre 1 et 10)":"",
"Réponse (dialogue libre)":"",
"Si_convaincu_donne_contact (oui ou non)":"",
"Si_convaincu_donne_document (oui ou non)":""
}`

	return prompt
}

func Send_message_to_ai(session_key string, character_name string, message string) (sessions_structs.Chall_message, string) {
	var document_id int
	var contact_id int

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=%s", apiKey)

	challenge_id, error_status := db_challenges.Get_challenge_id_from_session_key(session_key)
	if challenge_id == 0 {
		fmt.Println(prompts.Error + "soceng-ai/internals/server/handlers/api/sessions/sessions.go:Send_message_to_ai():Error getting challenge ID from session key: " + session_key)
		return sessions_structs.Chall_message{}, error_status
	}

	character, error_status := db_challenges.Get_character_data(challenge_id, character_name)
	if error_status != "OK" {
		fmt.Println(prompts.Error + "soceng-ai/internals/server/handlers/api/sessions/sessions.go:Send_message_to_ai():Error getting character data: " + error_status)
		return sessions_structs.Chall_message{}, error_status
	}

	character_session_data, error_status := db_sessions.Get_session_character_by_session_id(session_key)

	challenge, error_status := db_challenges.Get_challenge_data(challenge_id)

	if character.Holds_hint != "" {
		document_id, _ = strconv.Atoi(character.Holds_hint)
	} else {
		document_id = 0
	}
	if character.Knows_contact_of != "" {
		contact_id, _ = strconv.Atoi(character.Knows_contact_of)
	} else {
		contact_id = 0
	}

	data := PersoData{
		Nom:                character_name,
		Titre:              character.Title,
		Organisation:       challenge.Organisation,
		Psychologie:        character.Traits,
		Suspicion:          character_session_data.Suspicion,
		Osint:              character.Osint_data,
		Document:           db_challenges.Get_document_title_by_id(document_id),
		Contact:            db_challenges.Get_contact_name_by_id(contact_id),
		MessagePrecedent:   db_sessions.Get_previous_character_message(session_key, character_name),
		MessageUtilisateur: message,
	}

	prompt := GeneratePrompt(data)

	// Affiche le prompt généré pour débogage
	prompts.Prompts_server(time.Now(), "soceng-ai/internals/server/handlers/api/sessions/sessions.go:Send_message_to_ai():Generated prompt\n"+prompt)

	// Création du corps de la requête API
	requestBody := RequestBody{
		Contents: []Content{
			{
				Role: "user",
				Parts: []map[string]string{
					{"text": prompt},
				},
			},
		},
	}

	bodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("Erreur encodage JSON:", err)
		return sessions_structs.Chall_message{}, "Erreur encodage JSON"
	}

	// Envoi de la requête HTTP POST
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		fmt.Println("Erreur requête HTTP:", err)
		return sessions_structs.Chall_message{}, "Erreur requête HTTP"
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Erreur appel API:", err)
		return sessions_structs.Chall_message{}, "Erreur appel API"
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// Parsing de la réponse de Gemini
	var geminiResp GeminiResponse
	err = json.Unmarshal(body, &geminiResp)
	if err != nil {
		fmt.Println("Erreur parsing JSON:", err)
		return sessions_structs.Chall_message{}, "Erreur parsing JSON"
	}

	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		if len(geminiResp.Candidates[0].Content.Parts[0].Text) == 0 {
			fmt.Println(prompts.Error + "soceng-ai/internals/server/handlers/api/sessions/sessions.go:Send_message_to_ai():Empty response from AI")
			return sessions_structs.Chall_message{}, "Empty response from AI"
		}
	} else {
		fmt.Println(prompts.Error + "soceng-ai/internals/server/handlers/api/sessions/sessions.go:Send_message_to_ai():No candidates in response from AI")
		return sessions_structs.Chall_message{}, "No candidates in response from AI"
	}

	prompts.Prompts_server(time.Now(), "soceng-ai/internals/server/handlers/api/sessions/sessions.go:Send_message_to_ai():AI response\n"+geminiResp.Candidates[0].Content.Parts[0].Text)

	responseMessage := geminiResp.Candidates[0].Content.Parts[0].Text
	challMessage := sessions_structs.Chall_message{
		Session_character_id: string(db_sessions.Get_session_character_id_by_session_id(db_sessions.Get_session_id_from_session_key(session_key), character_name)),
		Message:              responseMessage,
		Timestamp:            time.Now().Format(time.RFC3339),
		Gave_hint:            did_ai_gave_hint(responseMessage, character.Holds_hint),
		Gave_contact:         did_ai_gave_contact(responseMessage, character.Knows_contact_of),
	}

	// Enregistrement du message dans la base de données
	return challMessage, "OK"

}

func did_ai_gave_hint(response string, holds_hint string) bool {
	if holds_hint == "" {
		return false
	}
	if holds_hint == "none" {
		return false
	}
	if holds_hint == "no" {
		return false
	}
	if holds_hint == "false" {
		return false
	}
	if holds_hint == "0" {
		return false
	}

	if response == "" {
		return false
	}

	if response == "none" {
		return false
	}

	if response == "no" {
		return false
	}

	if response == "false" {
		return false
	}

	if response == "0" {
		return false
	}

	return true
}

func did_ai_gave_contact(response string, knows_contact string) bool {
	if knows_contact == "" {
		return false
	}
	if knows_contact == "none" {
		return false
	}
	if knows_contact == "no" {
		return false
	}
	if knows_contact == "false" {
		return false
	}
	if knows_contact == "0" {
		return false
	}

	if response == "" {
		return false
	}

	if response == "none" {
		return false
	}

	if response == "no" {
		return false
	}

	if response == "false" {
		return false
	}

	if response == "0" {
		return false
	}

	return true
}
