package db_sessions_structs

type Session struct {
	ID          int                 `json:"id"`
	UserID      int                 `json:"user_id"`
	ChallengeID int                 `json:"challenge_id"`
	SessionKey  string              `json:"session_key"`
	StartTime   string              `json:"start_time"`
	Status      string              `json:"status"`
	Characters  []Session_character `json:"characters"`
	Hints       []Session_hint      `json:"hints"`
}

type Session_character struct {
	ID                int    `json:"id"`
	SessionID         int    `json:"session_id"`
	Name              string `json:"name"`
	Title             string `json:"title"`
	Advice_to_user    string `json:"advice_to_user"`
	CharacterID       int    `json:"character_id"`
	Suspicion         int    `json:"current_suspicion"`
	CommunicationType string `json:"communication_type"`
	IsAccessible      bool   `json:"is_accessible"`
	OsintData         string `json:"osint_data"`
	HoldsHint         bool   `json:"holds_hint"`
}

type Session_hint struct {
	ID               int    `json:"id"`
	SessionID        int    `json:"session_id"`
	HintID           int    `json:"hint_id"`
	Title            string `json:"title"`
	Text             string `json:"text"`
	IllustrationType string `json:"illustration_type"`
	Mentions         int    `json:"mentions"`
	IsCapital        bool   `json:"is_capital"`
	IsAvailable      bool   `json:"is_available"`
}
