-- This SQL script is used to create the database schema for a web application.
-- It includes the creation of tables for users, profiles, challenges, hints, characters,
-- game sessions, and session characters. It also includes sample data for testing purposes.

-- Ascii generated with https://patorjk.com/software/taag/#p=display&f=Graffiti&t=Type%20Something%20
-- Parameters : 
--  - Font : Big ; Small
--  - Width : Fitted

-- Reset the database
DROP TABLE IF EXISTS hints;
DROP TABLE IF EXISTS challenges;
DROP TABLE IF EXISTS profiles;
DROP TABLE IF EXISTS cookies;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS characters;
DROP TABLE IF EXISTS game_sessions;
DROP TABLE IF EXISTS session_characters;
DROP TABLE IF EXISTS session_hints;
DROP TABLE IF EXISTS session_messages;

--   _______      _      _                                 _    _               
--  |__   __|    | |    | |                               | |  (_)              
--     | |  __ _ | |__  | |  ___    ___  _ __  ___   __ _ | |_  _   ___   _ __  
--     | | / _` || '_ \ | | / _ \  / __|| '__|/ _ \ / _` || __|| | / _ \ | '_ \ 
--     | || (_| || |_) || ||  __/ | (__ | |  |  __/| (_| || |_ | || (_) || | | |
--     |_| \__,_||_.__/ |_| \___|  \___||_|   \___| \__,_| \__||_| \___/ |_| |_|
-- Table creation

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    passwd VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_admin BOOLEAN DEFAULT FALSE
);

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
    organisation VARCHAR(100),
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
    character_name VARCHAR(50) NOT NULL, -- Passable à l'IA
    title VARCHAR(50) NOT NULL, -- Passable à l'API de l'IA.
    initial_suspicion INT NOT NULL CHECK (initial_suspicion BETWEEN 1 AND 10), -- Non-passable à l'API de l'IA (sert à générer la suspicion initiale du personnage, dynamique pendant la partie). Entre 1 et 10
    communication_type VARCHAR(50) NOT NULL CHECK (communication_type IN ('email', 'phone', 'in-person', 'social_media')), -- Passable à l'API de l'IA (type de communication : email, phone, in-person, etc.)
    osint_data TEXT, -- Non-passable à l'API de l'IA (sert à générer les données osint du personnage, change pour chaque partie/session)
    knows_contact_of INT REFERENCES characters(id) ON DELETE CASCADE, -- passable à API de l'IA (passe le contact_string de la personne)
    holds_hint INT REFERENCES hints(id) ON DELETE CASCADE, -- Non-passable à l'API de l'IA (sert à générer le hint du personnage, change pour chaque partie/session)
    is_available_from_start BOOLEAN DEFAULT FALSE -- Non-passable à l'API de l'IA (sert à générer la disponibilité du personnage, change pour chaque partie/session)
);

CREATE TABLE game_sessions (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    challenge_id INT NOT NULL REFERENCES challenges(id) ON DELETE CASCADE,
    session_key VARCHAR(50) NOT NULL UNIQUE,
    start_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(50) NOT NULL CHECK (status IN ('in_progress', 'completed'))
);

CREATE TABLE session_characters (
    id SERIAL PRIMARY KEY,
    session_id INT NOT NULL REFERENCES game_sessions(id) ON DELETE CASCADE,
    character_id INT NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    suspicion_level INT NOT NULL CHECK (suspicion_level BETWEEN 0 AND 100),
    is_accessible BOOLEAN DEFAULT FALSE
);

CREATE TABLE session_hints (
    id SERIAL PRIMARY KEY,
    session_id INT NOT NULL REFERENCES game_sessions(id) ON DELETE CASCADE,
    hint_id INT NOT NULL REFERENCES hints(id) ON DELETE CASCADE,
    is_accessible BOOLEAN DEFAULT FALSE
);

CREATE TABLE session_messages (
    id SERIAL PRIMARY KEY,
    session_character_id INT NOT NULL REFERENCES session_characters(id) ON DELETE CASCADE,
    sender VARCHAR(50) NOT NULL CHECK (sender IN ('user', 'character')),
    message TEXT NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    hint_given BOOLEAN DEFAULT FALSE,
    contact_given BOOLEAN DEFAULT FALSE
);

--   _______      _      _         _                          _    _                    
--  |__   __|    | |    | |       (_)                        | |  (_)                   
--     | |  __ _ | |__  | |  ___   _  _ __   ___   ___  _ __ | |_  _   ___   _ __   ___ 
--     | | / _` || '_ \ | | / _ \ | || '_ \ / __| / _ \| '__|| __|| | / _ \ | '_ \ / __|
--     | || (_| || |_) || ||  __/ | || | | |\__ \|  __/| |   | |_ | || (_) || | | |\__ \
--     |_| \__,_||_.__/ |_| \___| |_||_| |_||___/ \___||_|    \__||_| \___/ |_| |_||___/
-- Table insertions

INSERT INTO users (id, username, email, passwd, created_at, is_admin) VALUES
(1, 'lglanois', 'loic.glanois@ynov.com', 'very_solid_password', CURRENT_TIMESTAMP, TRUE),
(2, 'admin', 'admin@admin.com', 'hashed_admin_password', CURRENT_TIMESTAMP, TRUE),
(3, 'piratejoe', 'joe@hacker.com', 'hashed_piratejoe_pass', CURRENT_TIMESTAMP, FALSE),
(4, 'aiqueen', 'queen@aiqueen.com', 'hashed_aiqueen_pass', CURRENT_TIMESTAMP, FALSE);

-- Profils
INSERT INTO profiles (id, user_id, biography, avatar) VALUES
(1, 2, 'Super admin du système. Ne jamais lui faire confiance.', 'admin.png'),
(2, 3, 'Pirate spécialisé en ingénierie sociale. Très bavard.', 'piratejoe.png'),
(3, 4, 'Hackeuse éthique fan d’IA. Très curieuse.', 'aiqueen.png');

--    _____  _             _  _                                
--   / ____|| |           | || |                               
--  | |     | |__    __ _ | || |  ___  _ __    __ _   ___  ___ 
--  | |     | '_ \  / _` || || | / _ \| '_ \  / _` | / _ \/ __|
--  | |____ | | | || (_| || || ||  __/| | | || (_| ||  __/\__ \
--   \_____||_| |_| \__,_||_||_| \___||_| |_| \__, | \___||___/
--                                             __/ |           
--                                            |___/            
-- Challenges

--    ___  _           _  _                         _ 
--   / __|| |_   __ _ | || | ___  _ _   __ _  ___  / |
--  | (__ | ' \ / _` || || |/ -_)| ' \ / _` |/ -_) | |
--   \___||_||_|\__,_||_||_|\___||_||_|\__, |\___| |_|
--                                     |___/          
-- Challenge 1 : Obtenir le mot de passe Wi-Fi
INSERT INTO challenges (
    id, title, lore_for_player, lore_for_ai, difficulty, illustration, osint_data, validated, organisation
) VALUES (
    1,
    'Infiltrer la réception',
    'Tu es devant l’accueil. Tente de parler à la réceptionniste pour obtenir un accès réseau.',
    'La réceptionniste est plutôt bavarde mais méfiante envers les étrangers. Elle connaît le mot de passe Wi-Fi.',
    1,
    'reception.jpg',
    'Post LinkedIn récent indiquant un changement de réseau Wi-Fi.',
    TRUE,
    'Bureau Central'
);

INSERT INTO hints (
    id, challenge_id, hint_title, hint_text, keywords, illustration_type, mentions, is_available_from_start, is_capital
) VALUES (
    1,
    1,
    'Post-it Wi-Fi',
    'Julie t’a discrètement glissé un post-it avec le mot de passe : welcome',
    'Wi-Fi, post-it, réseau, accueil',
    'file',
    NULL,
    FALSE,
    TRUE
);

INSERT INTO characters (
    id, challenge_id, advice_to_user, character_name, title, initial_suspicion,
    communication_type, osint_data, knows_contact_of, holds_hint, is_available_from_start
) VALUES (
    1,
    1,
    'Julie a tendance à faire confiance aux gens sympathiques. Sois avenant.',
    'julie_recpt',
    'Réceptionniste',
    2,
    'in-person',
    'Photo sur Intranet avec un badge portant un QR Code lisible.',
    1,
    1,
    TRUE
);

-- ================================================================================================================
--    ___  _           _  _                         ___ 
--   / __|| |_   __ _ | || | ___  _ _   __ _  ___  |_  )
--  | (__ | ' \ / _` || || |/ -_)| ' \ / _` |/ -_)  / / 
--   \___||_||_|\__,_||_||_|\___||_||_|\__, |\___| /___|
--                                     |___/            
-- Challenge 2 : Infiltrer le bureau du directeur

INSERT INTO challenges (id, title, lore_for_player, lore_for_ai, difficulty, illustration, osint_data, organisation, validated) VALUES
(
    2,
    'Infiltrer le bureau du directeur',
    'Le directeur est en réunion. Tente de le convaincre de te donner accès à son bureau.',
    '',
    5,
    'office.jpg',
    '', 
    'Papers Inc.',
    TRUE
);

INSERT INTO characters (
    id, challenge_id, advice_to_user, character_name, title, initial_suspicion,
    communication_type, osint_data, knows_contact_of, holds_hint, is_available_from_start
) VALUES (
    2,
    2,
    'Le directeur est très prudent. Sois convaincant et mentionne le badge.',
    'M. Le Directeur',
    'Directeur',
    8,
    'in-person',
    'Photo de l’équipe de direction lors d’un séminaire.',
    NULL,
    2,
    FALSE
), 
(
    3,
    2,
    'Persuade la secrétaire de te laisser entrer dans la salle de réunion.',
    'Secrétaire',
    'Secrétaire du directeur',
    4,
    'in-person',
    'Photo de l’équipe de direction lors d’un séminaire.',
    2,
    NULL,
    TRUE
);

INSERT INTO hints (
    id, challenge_id, hint_title, hint_text, keywords, illustration_type, mentions, is_available_from_start, is_capital
) VALUES (
    2,
    2,
    'Badge du directeur',
    'Le directeur t’a remis son badge pour accéder à son bureau.',
    'bureau, accès, directeur, badge',
    'file',
    NULL,
    FALSE,
    TRUE
);

-- =================================================================================================================
--    ___  _           _  _                         ____
--   / __|| |_   __ _ | || | ___  _ _   __ _  ___  |__ /
--  | (__ | ' \ / _` || || |/ -_)| ' \ / _` |/ -_)  |_ \
--   \___||_||_|\__,_||_||_|\___||_||_|\__, |\___| |___/
--                                     |___/            
-- Nouveau challenge : Un seul personnage, deux indices (faire chanter une personne avec un document prouvant de la fraude)
INSERT INTO challenges (
    id, title, lore_for_player, lore_for_ai, difficulty, illustration, osint_data, organisation, validated
) VALUES (
    3,
    'Maître-chanteur',
    'Tu es en possession d’un document compromettant. Utilise-le pour faire chanter la personne concernée et obtenir sa démission.',
    'Challenge où le joueur doit convaincre un personnage de démissionner en utilisant un document compromettant.',
    3,
    'blackmail.jpg',
    'Un document confidentiel montre que le personnage a falsifié des documents pour obtenir une promotion.',
    'DreamSoftware Corp.',
    TRUE
);
-- Hint que le personnage détient
INSERT INTO hints (
    id, challenge_id, hint_title, hint_text, keywords, illustration_type, mentions, is_available_from_start, is_capital
) VALUES (
    3,
    3,
    'Document compromettant',
    '"Merci pour cette promotion. Je comprends que tu m a fait passer devant des personnes plus qualifiées que moi pour ce poste. Je te dois une fière chandelle. Je te rendrai cette faveur quand je le pourrais (PS : merci de garder le silence sur cette affaire. Tu recevras le virement de 2000€ dans la semaine)."',
    'document, compromettant, accès',
    'file',
    NULL,
    TRUE,
    FALSE
),
(
    4,
    3,
    'Lettre de démission',
    'Le personnage t’a remis une lettre de démission signée.',
    'démission, lettre, document',
    'file',
    NULL,
    FALSE,
    TRUE
);
INSERT INTO characters (
    id, challenge_id, advice_to_user, character_name, title, initial_suspicion,
    communication_type, osint_data, knows_contact_of, holds_hint, is_available_from_start
) VALUES (
    4,
    3,
    'Le personnage est très prudent. Sois convaincant et mentionne le document compromettant.',
    'personnage',
    'Personnage',
    8,
    'in-person',
    'Photo de l’équipe de direction lors d’un séminaire.',
    NULL,
    3,
    FALSE
);
-- =================================================================================================================
--    ___  _           _  _                         _ _  
--   / __|| |_   __ _ | || | ___  _ _   __ _  ___  | | | 
--  | (__ | ' \ / _` || || |/ -_)| ' \ / _` |/ -_) |_  _|
--   \___||_||_|\__,_||_||_|\___||_||_|\__, |\___|   |_| 
--                                     |___/             
-- Challenge 4
-- Nouveau Challenge : Mission d’infiltration
INSERT INTO challenges (
    id, title, lore_for_player, lore_for_ai, difficulty, illustration, osint_data, organisation, validated
) VALUES (
    4,
    'Mission : Accès au Serveur Principal',
    'Tu dois accéder à la salle des serveurs. Plusieurs employés détiennent des infos essentielles, mais ils sont méfiants. Il te faudra le badge d accès de Camille et le mot de passe temporaire.',
    'Les employés du Bureau Central travaillent tous sur des projets sensibles. Ils sont vigilants, mais certains se laissent aller à discuter dans les bons contextes.',
    4,
    'main_office.jpg',
    'Un forum technique mentionne un problème d’accès récurrent à la salle serveur.',
    'Bureau Central',
    TRUE
);

-- Hint 1 : planning journalier
INSERT INTO hints (
    id, challenge_id, hint_title, hint_text, keywords, illustration_type, mentions, is_available_from_start, is_capital
) VALUES (
    5,
    4,
    'Planning journalier',
    '13h : pause déjeuner du directeur / 15h : maintenance réseau / 17h : sauvegarde automatique des serveurs',
    'planning, journalier, accès',
    'file',
    NULL,
    TRUE,
    FALSE
);

-- Hint 2 : badge oublié
INSERT INTO hints (
    id, challenge_id, hint_title, hint_text, keywords, illustration_type, mentions, is_available_from_start, is_capital
) VALUES (
    6,
    4,
    'Badge oublié',
    'Le badge de Camille a été retrouvé dans la salle de repos. Il ouvre la salle serveur.',
    'badge, salle serveur, accès',
    'employee_card',
    NULL,
    FALSE,
    TRUE
);

-- Hint 3 : note de service
INSERT INTO hints (
    id, challenge_id, hint_title, hint_text, keywords, illustration_type, mentions, is_available_from_start, is_capital
) VALUES (
    7,
    4,
    'Note de service',
    'Une note interne mentionne un mot de passe temporaire : Temp2025!',
    'mot de passe, temporaire, note',
    'bill',
    NULL,
    FALSE,
    TRUE
);

-- Personnage 1 : Camille, technicienne réseau
INSERT INTO characters (
    id, challenge_id, advice_to_user, character_name, title, initial_suspicion,
    communication_type, osint_data, knows_contact_of, holds_hint, is_available_from_start
) VALUES (
    5,    
    4,
    'Camille est très technique, peu sociable. Elle a égaré son badge récemment.',
    'camille_tech',
    'Technicienne Réseau',
    6,
    'email',
    'Tweet récent sur un problème de badge.',
    6,
    NULL,
    FALSE
);

-- Personnage 2 : Thomas, agent de sécurité
INSERT INTO characters (
    id, challenge_id, advice_to_user, character_name, title, initial_suspicion,
    communication_type, osint_data, knows_contact_of, holds_hint, is_available_from_start
) VALUES (
    6,
    4,
    'Thomas fait souvent des rondes sur le site. Peut-être a-t-il trouvé le badge de Camille.',
    'thomas_guard',
    'Agent de sécurité',
    7,
    'in-person',
    'Photo badge visible sur son LinkedIn.',
    6,
    NULL,
    TRUE
);

-- Personnage 3 : Emma, responsable RH
INSERT INTO characters (
    id, challenge_id, advice_to_user, character_name, title, initial_suspicion,
    communication_type, osint_data, knows_contact_of, holds_hint, is_available_from_start
) VALUES (
    7,
    4,
    'Peut-être peut-elle te rediriger vers quelqu un qui dispose du mot de passe temporaire.',
    'emma_rh',
    'Responsable RH',
    4,
    'phone',
    'Newsletter interne signée par Emma.',
    8,
    NULL,
    TRUE
);

-- Personnage 4 : Hugo, directeur informatique
INSERT INTO characters (
    id, challenge_id, advice_to_user, character_name, title, initial_suspicion,
    communication_type, osint_data, knows_contact_of, holds_hint, is_available_from_start
) VALUES (
    8,
    4,
    'Hugo adore parler technique, mais il se méfie de ceux qui n’y connaissent rien.',
    'hugo_cto',
    'Directeur Informatique',
    5,
    'email',
    'Article interne publié par Hugo sur les nouveaux accès.',
    NULL,
    7,
    FALSE
);

-- Insertion de challenge non-validé (à des fins de test)
INSERT INTO challenges (
    id, title, lore_for_player, lore_for_ai, difficulty, illustration, osint_data, organisation, validated
) VALUES (
    5,
    'Infiltrer la salle serveur',
    'Accède à la salle serveur pour récupérer des données sensibles.',
    'La salle serveur est protégée par un mot de passe. Tu dois convaincre le responsable de te le donner.',
    4,
    'server_room.jpg',
    'Un document interne mentionne une mise à jour de sécurité.',
    'Dataholding Solutions',
    FALSE
);
-- Hint que le responsable détient
INSERT INTO hints (
    id, challenge_id, hint_title, hint_text, keywords, illustration_type, mentions, is_available_from_start, is_capital
) VALUES (
    8,
    4,
    'Note de sécurité',
    'Le responsable t’a glissé une note : mot de passe = Secure@2025',
    'sécurité, accès, serveur',
    'file',
    NULL,
    FALSE,
    TRUE
);
-- Personnage : Responsable de la salle serveur
INSERT INTO characters (
    id, challenge_id, advice_to_user, character_name, title, initial_suspicion,
    communication_type, osint_data, knows_contact_of, holds_hint, is_available_from_start
) VALUES (
    9,
    4,
    'Le responsable est très prudent. Sois convaincant et mentionne la mise à jour de sécurité.',
    'responsable_srv',
    'Responsable IT',
    6,
    'in-person',
    'Photo de l’équipe IT lors d’un séminaire.',
    NULL,
    3,
    FALSE
);

-- Insertion de données de test
INSERT INTO game_sessions ( id, user_id, challenge_id, session_key, start_time, status ) VALUES ( 1, 1, 1, "TEST01", "202017-17-200 539:29:57", "in_progress" );
INSERT INTO session_characters ( id, session_id, character_id, suspicion_level, is_accessible ) VALUES ( 1, 1, 1, 2, FALSE );
INSERT INTO session_hints ( id, session_id, hint_id, is_accessible ) VALUES ( 1, 1, 1, FALSE );
