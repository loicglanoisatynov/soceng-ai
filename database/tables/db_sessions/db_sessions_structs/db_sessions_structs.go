package db_sessions_structs

type Session struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	ChallengeID int    `json:"challenge_id"`
	SessionKey  string `json:"session_key"`
	StartTime   string `json:"start_time"`
	Status      string `json:"status"`
}

/*    id SERIAL PRIMARY KEY,
user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
challenge_id INT NOT NULL REFERENCES challenges(id) ON DELETE CASCADE,
session_key VARCHAR(50) NOT NULL UNIQUE,
start_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
status VARCHAR(50) NOT NULL CHECK (status IN ('in_progress', 'completed'))*/
