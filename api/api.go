package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	prompt := `Ce message vient d'une API de CTF d'ingénierie sociale, tu dois répondre en fonction des données de jeu passées. Tu as le droit de modifier les données identifiées Altérables, tu as le droit de toucher aux autres. Réponds conformément à ton personnage dans la zone de texte Réponse.
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

func main() {
	// Clé API
	apiKey := "AIzaSyDVLqTrPlzSI3KvEJHEN58uFQeRPY3rPTU"
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=%s", apiKey)

	// 👤 Données du personnage (exemple)
	data := PersoData{
		Nom:                "Jean Dupont",
		Titre:              "Technicien réseau",
		Organisation:       "ZetaCorp",
		Psychologie:        "naïf, serviable, stressé",
		Suspicion:          4,
		Osint:              "utilise gmail, a un chat nommé Minou",
		Document:           "schema_reseau.pdf",
		Contact:            "alice.it@zetacorp.com",
		MessagePrecedent:   "Je vous avais dit que je devais demander à mon supérieur...",
		MessageUtilisateur: "Tu peux me l'envoyer maintenant ? C'est urgent.",
	}

	prompt := GeneratePrompt(data)

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
		return
	}

	// Envoi de la requête HTTP POST
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		fmt.Println("Erreur requête HTTP:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Erreur appel API:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// Affichage brut pour debug
	fmt.Println("Réponse brute :")
	fmt.Println(string(body))

	// Parsing de la réponse de Gemini
	var geminiResp GeminiResponse
	err = json.Unmarshal(body, &geminiResp)
	if err != nil {
		fmt.Println("Erreur parsing JSON:", err)
		return
	}

	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		fmt.Println("\n--- Réponse du personnage ---")
		fmt.Println(geminiResp.Candidates[0].Content.Parts[0].Text)
	} else {
		fmt.Println("Réponse vide ou mal formée.")
	}
}
