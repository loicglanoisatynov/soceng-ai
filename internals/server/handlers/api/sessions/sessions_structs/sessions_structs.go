package sessions_structs

type Create_session_request struct {
	Challenge_name string `json:"title"`
}

type Get_challenge_response struct {
	Challenge_name         string            `json:"chall_name"`
	Challenge_desc         string            `json:"chall_desc"`
	Challenge_illustration string            `json:"illustration"`
	Session_key            string            `json:"session_key"`
	Chall_characters       []Chall_character `json:"chall_characters"`
	Chall_hints            []Chall_hint      `json:"chall_hints"`
	Chall_status           string            `json:"chall_status"`
}

type Chall_character struct {
	Name               string `json:"name"`
	Title              string `json:"title"`
	Advice_to_user     string `json:"advice_to_user"`
	Suspicion          int    `json:"suspicion"`
	Communication_type string `json:"communication_type"`
	Osint_data         string `json:"osint_data"`
}

type Chall_hint struct {
	Title             string `json:"title"`
	Text              string `json:"text"`
	Illustration_type string `json:"illustration_type"`
	Mentions          int    `json:"mentions"`
	Is_capital        bool   `json:"is_capital"`
}

type Chall_message struct {
	User_or_character    string `json:"user_or_character"`
	Message              string `json:"message"`
	Timestamp            string `json:"timestamp"`
	Gave_hint            bool   `json:"gave_hint"`
	Gave_contact         bool   `json:"gave_contact"`
	Session_character_id string `json:"session_character_id"`
}

type Post_session_data_request struct {
	Character_name string `json:"character_name"`
	Message        string `json:"message"`
}

// L'autre sous-route de /api/sessions est la sous-route des clés de session (**/api/sessions/EXMPLE**) qui récupère toutes les données de session du challenge en question. Ces données sont : le nom du challenge (`json:"chall_name"`), la description du challenge (`json:"chall_desc"`), le nom de l'image d'illustration du challenge (`json:"illustration"`), la clé de session du challenge (pour qu'elle ne disparaisse pas ; `json:"session_key"`), un **array** contenant tous les personnages de la session (`json:"chall_characters"`) qui sont d'autres objets dont je vais t'énumérer également les attributs plus bas, les hints de la session (`json:"chall_hints"`) même chose que pour les personnages et le status du challenge (`json:"chall_status"`). Je précise que lorsque la mécanique des messages sera ajoutée, ils seront également dans ce bloc de json et posséderont leur propre array
// Les personnages de la session sont des objets qui possèdent les attributs suivants : le nom du personnage (`json:"name"`), le titre du personnage (`json:"title"`), le conseil à l'utilisateur (`json:"advice_to_user"`), le niveau de suspicion (`json:"suspicion"`), le type de communication (`json:"communication_type"`) et les données OSINT (`json:"osint_data"`). Les hints de la session sont des objets qui possèdent les attributs suivants : le titre du hint (`json:"title"`), le texte du hint (`json:"text"`), le type d'illustration du hint (`json:"illustration_type"`), le nombre de mentions du hint dans la session (`json:"mentions"`), un booléen qui indique si le hint est capital ou non (`json:"is_capital"`) et un booléen qui indique si le hint est visible ou non (`json:"is_visible"`)
// Les message du challenge auront les attributs `json:"user_or_character"`, `json:"message"`, `json:"character_name"`, `json:"timestamp"`, `json:"holds_hint"`
