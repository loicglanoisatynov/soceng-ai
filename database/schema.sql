DROP TABLE IF EXISTS hints;
DROP TABLE IF EXISTS challenges;
DROP TABLE IF EXISTS profiles;
DROP TABLE IF EXISTS cookies;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS characters;

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

-- Challenge : Obtenir le mot de passe Wi-Fi
INSERT INTO challenges (
    id, title, lore_for_player, lore_for_ai, difficulty, illustration, osint_data
) VALUES (
    1,
    'Infiltrer la réception',
    'Tu es devant l’accueil. Tente de parler à la réceptionniste pour obtenir un accès réseau.',
    'La réceptionniste est plutôt bavarde mais méfiante envers les étrangers. Elle connaît le mot de passe Wi-Fi.',
    1,
    'reception.jpg',
    'Post LinkedIn récent indiquant un changement de réseau Wi-Fi.'
);

-- Hint : Récompense du challenge
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

-- Personnage : Julie la réceptionniste
INSERT INTO characters (
    id, challenge_id, advice_to_user, symbolic_name, title, initial_suspicion,
    communication_type, symbolic_osint_data, knows_contact_of, holds_hint, is_available_from_start
) VALUES (
    1,
    1,
    'Julie a tendance à faire confiance aux gens sympathiques. Sois avenant.',
    'julie_recpt',
    'Réceptionniste',
    20,
    'in-person',
    'Photo sur Intranet avec un badge portant un QR Code lisible.',
    1,
    1,
    TRUE
);

-- Nouveau challenge : Convaincre deux employés
INSERT INTO challenges (
    id, title, lore_for_player, lore_for_ai, difficulty, illustration, osint_data
) VALUES (
    2,
    'Récupérer les accès internes',
    'Deux employés possèdent chacun une moitié d’une information précieuse. Obtiens leur confiance.',
    'Le premier personnage (Paul) peut te rediriger vers sa collègue (Claire) qui a le complément. Ils sont prudents, mais pas impossibles à convaincre.',
    3,
    'office_access.jpg',
    'Un document de réunion interne montre que Paul et Claire travaillent sur le même projet.'
);

-- Hint que Claire détient
INSERT INTO hints (
    id, challenge_id, hint_title, hint_text, keywords, illustration_type, mentions, is_available_from_start, is_capital
) VALUES (
    2,
    2,
    'Mémo technique',
    'Claire t’a transmis un mémo confidentiel : mot de passe = Internal@2025',
    'accès, interne, projet',
    'file',
    NULL,
    FALSE,
    TRUE
);

-- Personnage 1 : Paul (oriente vers Claire)
INSERT INTO characters (
    id, challenge_id, advice_to_user, symbolic_name, title, initial_suspicion,
    communication_type, symbolic_osint_data, knows_contact_of, holds_hint, is_available_from_start
) VALUES (
    2,
    2,
    'Paul est méthodique. Il ne donne rien sans preuve, mais il t’orientera si tu sembles bien renseigné.',
    'paul_dev',
    'Développeur',
    40,
    'email',
    'Paul est actif sur GitHub, souvent la nuit.',
    3,
    NULL,
    TRUE
);

-- Personnage 2 : Claire (détient le hint)
INSERT INTO characters (
    id, challenge_id, advice_to_user, symbolic_name, title, initial_suspicion,
    communication_type, symbolic_osint_data, knows_contact_of, holds_hint, is_available_from_start
) VALUES (
    3,
    2,
    'Claire est méfiante mais bavarde si tu mentionnes Paul et leur projet commun.',
    'claire_hr',
    'Chargée RH',
    50,
    'in-person',
    'Photo d’équipe avec Paul lors d’un team-building.',
    2,
    2,
    FALSE
);

-- Nouveau Challenge : Mission d’infiltration
INSERT INTO challenges (
    id, title, lore_for_player, lore_for_ai, difficulty, illustration, osint_data
) VALUES (
    3,
    'Mission : Accès au Serveur Principal',
    'Tu dois accéder à la salle des serveurs. Plusieurs employés détiennent des infos essentielles, mais ils sont méfiants.',
    'Les employés du Bureau Central travaillent tous sur des projets sensibles. Ils sont vigilants, mais certains se laissent aller à discuter dans les bons contextes.',
    4,
    'main_office.jpg',
    'Un forum technique mentionne un problème d’accès récurrent à la salle serveur.'
);

-- Hint 1 : planning journalier
INSERT INTO hints (
    id, challenge_id, hint_title, hint_text, keywords, illustration_type, mentions, is_available_from_start, is_capital
) VALUES (
    3,
    3,
    'Planning journalier',
    '13h : pause déjeuner du directeur / 15h : maintenance réseau / 17h : sauvegarde automatique des serveurs',
    'planning, horaire, serveur',
    'file',
    NULL,
    TRUE,
    FALSE
);

-- Hint 2 : badge oublié
INSERT INTO hints (
    id, challenge_id, hint_title, hint_text, keywords, illustration_type, mentions, is_available_from_start, is_capital
) VALUES (
    4,
    3,
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
    5,
    3,
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
    id, challenge_id, advice_to_user, symbolic_name, title, initial_suspicion,
    communication_type, symbolic_osint_data, knows_contact_of, holds_hint, is_available_from_start
) VALUES (
    4,
    3,
    'Camille est très technique, peu sociable. Elle a égaré son badge récemment.',
    'camille_tech',
    'Technicienne Réseau',
    60,
    'email',
    'Tweet récent sur un problème de badge.',
    5,
    4,
    FALSE
);

-- Personnage 2 : Thomas, agent de sécurité
INSERT INTO characters (
    id, challenge_id, advice_to_user, symbolic_name, title, initial_suspicion,
    communication_type, symbolic_osint_data, knows_contact_of, holds_hint, is_available_from_start
) VALUES (
    5,
    3,
    'Thomas ne plaisante pas. Il faut avoir une bonne couverture pour qu’il parle.',
    'thomas_guard',
    'Agent de sécurité',
    70,
    'in-person',
    'Photo badge visible sur son LinkedIn.',
    6,
    NULL,
    TRUE
);

-- Personnage 3 : Emma, responsable RH
INSERT INTO characters (
    id, challenge_id, advice_to_user, symbolic_name, title, initial_suspicion,
    communication_type, symbolic_osint_data, knows_contact_of, holds_hint, is_available_from_start
) VALUES (
    6,
    3,
    'Emma est sympathique, surtout si tu mentionnes Thomas ou les procédures internes.',
    'emma_rh',
    'Responsable RH',
    40,
    'phone',
    'Newsletter interne signée par Emma.',
    7,
    NULL,
    FALSE
);

-- Personnage 4 : Hugo, directeur informatique
INSERT INTO characters (
    id, challenge_id, advice_to_user, symbolic_name, title, initial_suspicion,
    communication_type, symbolic_osint_data, knows_contact_of, holds_hint, is_available_from_start
) VALUES (
    7,
    3,
    'Hugo adore parler technique, mais il se méfie de ceux qui n’y connaissent rien.',
    'hugo_cto',
    'Directeur Informatique',
    50,
    'email',
    'Article interne publié par Hugo sur les nouveaux accès.',
    8,
    5,
    FALSE
);

-- Personnage 5 : Léa, assistante administrative
INSERT INTO characters (
    id, challenge_id, advice_to_user, symbolic_name, title, initial_suspicion,
    communication_type, symbolic_osint_data, knows_contact_of, holds_hint, is_available_from_start
) VALUES (
    8,
    3,
    'Léa voit tout, entend tout. Elle partage parfois des infos sans s’en rendre compte.',
    'lea_admin',
    'Assistante Administrative',
    30,
    'in-person',
    'Cliché lors d’un pot d’entreprise où Camille et Léa discutent.',
    4,
    3,
    TRUE
);
