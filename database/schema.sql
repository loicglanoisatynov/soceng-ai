DROP TABLE IF EXISTS cookies;
DROP TABLE IF EXISTS profiles;
DROP TABLE IF EXISTS challenges;
DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    passwd VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_admin BOOLEAN DEFAULT FALSE
);

INSERT INTO users (id, username, email, passwd, created_at, is_admin) VALUES
(1, 'lglanois', 'loic.glanois@ynov.com', 'very_solid_password', CURRENT_TIMESTAMP, TRUE);

CREATE TABLE cookies (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    cookie_value VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_access TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE profiles (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    biography VARCHAR(255),
    avatar VARCHAR(255)
);

CREATE TABLE challenges (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    lore_for_player TEXT NOT NULL,
    lore_for_ai TEXT NOT NULL,
    difficulty INT NOT NULL CHECK (difficulty BETWEEN 1 AND 5),
    illustration VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    validated BOOLEAN DEFAULT FALSE,
    osint_data TEXT
);

CREATE TABLE hints (
    id SERIAL PRIMARY KEY,
    challenge_id INT NOT NULL REFERENCES challenges(id) ON DELETE CASCADE,
    hint_title VARCHAR(100) NOT NULL,
    hint_text TEXT NOT NULL,
    keywords TEXT,
    illustration_type VARCHAR(50) NOT NULL CHECK (illustration_type IN ('bill', 'employee_card', 'file')),
    mentions INT REFERENCES characters(id) ON DELETE CASCADE,
    is_available_from_start BOOLEAN DEFAULT FALSE,
    is_capital BOOLEAN DEFAULT FALSE
);

CREATE TABLE characters (
    id SERIAL PRIMARY KEY,
    challenge_id INT NOT NULL REFERENCES challenges(id) ON DELETE CASCADE,
    advice_to_user TEXT, -- Passable à l'API de l'IA.
    symbolic_name VARCHAR(50) NOT NULL, -- Non-passable à l'IA. Génère un nom de personnage aléatoire pour chaque partie.
    title VARCHAR(50) NOT NULL, -- Passable à l'API de l'IA.
    initial_suspicion INT NOT NULL, -- Non-passable à l'API de l'IA (sert à générer la suspicion initiale du personnage, dynamique pendant la partie).
    communication_type VARCHAR(50) NOT NULL CHECK (communication_type IN ('email', 'phone', 'in-person', 'social_media')), -- Passable à l'API de l'IA (type de communication : email, phone, in-person, etc.)
    symbolic_osint_data TEXT, -- Non-passable à l'API de l'IA (sert à générer les données osint du personnage, change pour chaque partie/session)
    knows_contact_of INT NOT NULL REFERENCES characters(id) ON DELETE CASCADE, -- passable à API de l'IA (passe le contact_string de la personne)
    holds_hint INT REFERENCES hints(id) ON DELETE CASCADE, -- Non-passable à l'API de l'IA (sert à générer le hint du personnage, change pour chaque partie/session)
    is_available_from_start BOOLEAN DEFAULT FALSE -- Non-passable à l'API de l'IA (sert à générer la disponibilité du personnage, change pour chaque partie/session)
);

-- Tables à créer : 
-- sessions (id, user_id, challenge_id, start_time, end_time, status)
-- game_characters (id, character_id, character_name, suspicion_level, is_contacted, is_suspect, session_id)

-- Utilisateurs
INSERT INTO users (id, username, email, passwd, created_at, is_admin) VALUES
(2, 'admin', 'admin@admin.com', 'hashed_admin_password', CURRENT_TIMESTAMP, TRUE),
(3, 'piratejoe', 'joe@hacker.com', 'hashed_piratejoe_pass', CURRENT_TIMESTAMP, FALSE),
(4, 'aiqueen', 'queen@aiqueen.com', 'hashed_aiqueen_pass', CURRENT_TIMESTAMP, FALSE);

-- Profils
INSERT INTO profiles (id, user_id, biography, avatar) VALUES
(1, 2, 'Super admin du système. Ne jamais lui faire confiance.', 'admin.png'),
(2, 3, 'Pirate spécialisé en ingénierie sociale. Très bavard.', 'piratejoe.png'),
(3, 4, 'Hackeuse éthique fan d’IA. Très curieuse.', 'aiqueen.png');
 
-- Challenges du jeu
-- INSERT INTO challenges (title, description, flag, difficulty) VALUES
-- (
--   'Infiltrer la réception',
--   'Discute avec la réceptionniste pour obtenir le mot de passe Wi-Fi.',
--   'FLAG{wifi}',
--   'Facile'
-- ),
-- (
--   'Convaincre le directeur',
--   'Tente de récupérer des infos techniques sans éveiller ses soupçons.',
--   'FLAG{tech}',
--   'Difficile'
-- ),
-- (
--   'Nettoyage stratégique',
--   'La femme de ménage en sait plus que tu ne crois. Profite de son bavardage.',
--   'FLAG{balaie}',
--   'Moyen'
-- );