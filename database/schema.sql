DROP TABLE IF EXISTS cookies;
DROP TABLE IF EXISTS profiles;
DROP TABLE IF EXISTS challenges;
DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    passwd VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
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
    description TEXT NOT NULL,
    flag VARCHAR(255) NOT NULL,
    difficulty VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);



-- Utilisateurs
INSERT INTO users (username, email, passwd) VALUES
('admin', 'admin@admin.com', 'hashed_admin_password'),
('piratejoe', 'joe@hacker.com', 'hashed_piratejoe_pass'),
('aiqueen', 'queen@aiqueen.com', 'hashed_aiqueen_pass');

-- Profils
INSERT INTO profiles (user_id, biography, avatar) VALUES
(1, 'Super admin du système. Ne jamais lui faire confiance.', 'admin.png'),
(2, 'Pirate spécialisé en ingénierie sociale. Très bavard.', 'piratejoe.png'),
(3, 'Hackeuse éthique fan d’IA. Très curieuse.', 'aiqueen.png');
 
-- Challenges du jeu
INSERT INTO challenges (title, description, flag, difficulty) VALUES
(
  'Infiltrer la réception',
  'Discute avec la réceptionniste pour obtenir le mot de passe Wi-Fi.',
  'FLAG{wifi}',
  'Facile'
),
(
  'Convaincre le directeur',
  'Tente de récupérer des infos techniques sans éveiller ses soupçons.',
  'FLAG{tech}',
  'Difficile'
),
(
  'Nettoyage stratégique',
  'La femme de ménage en sait plus que tu ne crois. Profite de son bavardage.',
  'FLAG{balaie}',
  'Moyen'
);
