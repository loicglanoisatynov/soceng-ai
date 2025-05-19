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
	User_or_character string `json:"user_or_character"`
	Message           string `json:"message"`
	Timestamp         string `json:"timestamp"`
	Holds_hint        bool   `json:"holds_hint"`
}
