const { Command } = require('commander');
const readline = require('readline').createInterface({
  input: process.stdin,
  output: process.stdout
});
const axios = require('axios');
const { start } = require('repl');
// const { read } = require('fs');

let socengai_username_cookie_value = '';
let socengai_auth_cookie_value = '';

// Fonction principale
function main() {
  let host;
  let port;

  async function initialise() {
    while (true) {
      let h = await new Promise((resolve) => {
        readline.question('Entrez le nom d\'hôte de l\'application distante (par défaut : localhost) : ', (input) => {
          resolve(input.trim() || 'localhost');
        });
      });

      let p = await new Promise((resolve) => {
        readline.question('Entrez le port de l\'application distante (par défaut : 80) : ', (input) => {
          resolve(input.trim() || '80');
        });
      });

      host = h;
      port = p;

      console.log(`Tentative de connexion à ${host}:${port}`);

      try {
        const url = `http://${host}:${port}/api/hello-world`;
        const response = await axios.get(url);

        if (response.data === "Hello, World !") {
          console.log('Serveur distant fonctionnel !');
          break;
        } else {
          console.log('Réponse inattendue du serveur. Veuillez réessayer.');
        }
      } catch (error) {
        console.error(`Erreur lors de la connexion au serveur : ${error.message}`);
      }
    }
  }

    async function loginMenu() {
      console.log('Menu :');
      console.log('1. Se connecter');
      console.log('2. Créer un compte');
      console.log('3. Quitter');

      let choix = await new Promise((resolve) => {
        readline.question('Entrez votre choix : ', (input) => {
          resolve(input.trim());
        });
      });

      switch (choix) {
        case '1':
          await seConnecter();
          break;
        case '2':
          await creerCompte();
          break;
        case '3':
          console.log('Au revoir !');
          readline.close();
          process.exit(0);
        default:
          console.log('Choix invalide. Veuillez réessayer.');
          await loginMenu();
          break;
      }
  }


  async function seConnecter() {
    let cookies_exist = false;
    let socengai_username_cookie_value = '';
    username_cookie = await check_cookies();
    if (username_cookie !== '') {
      console.log('Session existante détectée.');
      let reponse = await new Promise((resolve) => {
        readline.question('Une session est déjà active, au nom de : ' + username_cookie + '. Voulez-vous utiliser cette session ? (O/N)', (input) => {
          resolve(input.trim().toUpperCase());
        });
      });
      if (reponse === 'O' || reponse === 'OUI' || reponse === 'YES' || reponse === 'Y' || reponse === '') {
        const cookies_confirmed = await ask_for_cookies_confirmation();
        if (!cookies_confirmed) {
          console.log('Session disparue. Veuillez vous reconnecter.');
          await delete_cookies();
          await seConnecter();
        } else {
          console.log('Connexion réussie avec la session existante.');
          await dashboard_menu();
        }
      } else if (reponse === 'N') {
        console.log('Déconnexion de la session existante.');
        socengai_username_cookie_value = '';
        socengai_auth_cookie_value = '';
      } else {
        console.log('Choix invalide. Veuillez réessayer.');
        return;
      }
    }

    // Si aucune session n'est active, demande les identifiants
    const username = await new Promise((resolve) => {
      readline.question('Entrez votre nom d\'utilisateur : ', resolve);
    });
    const password = await new Promise((resolve) => {
      readline.question('Entrez votre mot de passe : ', resolve);
    });
    try {
      const url = `http://${host}:${port}/login`;
      const response = await axios.post(url, 
        JSON.stringify({
          username: username,
          password: password
        }), 
        {
          headers: {
            'Content-Type': 'application/json'
          }
        }
      );

      if (response.status === 200) {
        console.log('Connexion réussie !');
        await save_cookies(response.headers['set-cookie']);

        await dashboard_menu();
      } else {
        console.log('Erreur lors de la connexion :', response.data);
      }
    } catch (error) {
      if (error.response && error.response.data) {
        if (error.response.data === "User not found\n") {
          console.log('Utilisateur non trouvé');
        } else if (error.response.data === "Invalid password\n") {
          console.log('Mot de passe invalide');
        } else {
          console.log(error.response.data);
        }
      } else {
        console.log(`Erreur de connexion : ${error.message}`);
      }
    }
  }

  async function save_cookies(cookie) {
    socengai_username_cookie_value = cookie.find(c => c.startsWith('socengai-username='));
    if (!socengai_username_cookie_value) {
      console.error('Aucun cookie de nom d\'utilisateur trouvé.');
      return;
    }
    socengai_auth_cookie_value = cookie.find(c => c.startsWith('socengai-auth='));
    if (!socengai_auth_cookie_value) {
      console.error('Aucun cookie d\'authentification trouvé.');
      return;
    }
    
    const fs = require('fs').promises; // Utilisez la version promise de fs
    
    try {
        await fs.writeFile('cookies.txt', socengai_username_cookie_value);
        console.log('Cookies enregistrés avec succès dans cookies.txt');
        
        await fs.appendFile('cookies.txt', `\n${socengai_auth_cookie_value}`);
        console.log('Cookie d\'authentification ajouté avec succès dans cookies.txt');
    } catch (err) {
        console.error(`Erreur lors de l'enregistrement des cookies : ${err.message}`);
    }
  }

  async function menuChallenges() {
    while (true) {
      console.log('Menu des challenges :');
      console.log('1. Choisir un challenge');
      console.log('2. Quitter le menu des challenges');
      let choix = await new Promise((resolve) => {
        readline.question('Entrez votre choix : ', (input) => {
          resolve(input.trim());
        });
      });
      switch (choix) {
        case '1':
          await lister_challenges();
          break;
        case '2':
          return; // Retour au menu du tableau de bord
        default:
          console.log('Choix invalide. Veuillez réessayer.');
      }
    }
  }

  async function lister_challenges() {
    const url = `http://${host}:${port}/api/dashboard`;
    let challenges = [];
    try {
      const response = await axios.get(url, {
        headers: {
          'Cookie': `${socengai_username_cookie_value}; ${socengai_auth_cookie_value}`
        }
      });
      challenges = response.data.challenges;
      if (challenges.length === 0) {
        console.log('Aucun challenge disponible pour le moment.');
      } else {
        console.log('Challenges disponibles :');
        challenges.forEach((challenge, index) => {
          if (!challenge.session_key) {
            session_key = 'Aucune session';
          } else {
            session_key = challenge.session_key;
          }
          console.log(`${index + 1}. ${challenge.name} - ${challenge.description} (session : ${session_key})`);
        });
      }
      console.log();
    } catch (error) {
      console.error(`Erreur lors de la récupération des challenges : ${error.message}`);
    }

    // Demande à l'utilisateur de choisir un challenge
    const choix = await new Promise((resolve) => {
      readline.question('Entrez le numéro du challenge à rejoindre : ', (input) => {
        resolve(input.trim());
      });
    });
    const challengeIndex = parseInt(choix) - 1;
    if (isNaN(challengeIndex) || challengeIndex < 0 || challengeIndex >= challenges.length) {
      console.log('Choix invalide. Veuillez réessayer.');
      return;
    }
    const challenge = challenges[challengeIndex];
    console.log(`Vous avez rejoint le challenge : ${challenge.name}`);
    // Si le challenge a une session_key, demander si il veut recommencer la session ou reprendre la session
    if (challenge.session_key) {
      const reponse = await new Promise((resolve) => {
        readline.question('Voulez-vous reprendre la session existante (1) ou en créer une nouvelle ? (2) : ', (input) => {
          resolve(input.trim().toUpperCase());
        });
      });
      if (reponse === '1') {
        play_challenge(challenge);
      } else if (reponse === '2') {
        console.log(`Création d'une nouvelle session pour le challenge : ${challenge.name}`);
        start_challenge(challenge);
      } else {
        console.log('Choix invalide. Veuillez réessayer.');
      }
    } else {
      console.log(`Aucune session existante pour le challenge : ${challenge.name}`);
      start_challenge(challenge);
    }
  }

  async function play_challenge(challenge) {}

  async function dashboard_menu() {
    while (true) {
      let choix = 0;
      console.log('Menu du tableau de bord :');
      console.log('1. Afficher les informations de l\'utilisateur');
      console.log('2. Modifier les informations de l\'utilisateur');
      console.log('3. Menu challenges');
      console.log('4. Déconnexion');
      console.log('5. Quitter');
      
      choix = await new Promise((resolve) => {
        readline.question('Entrez votre choix : ', (input) => {
          resolve(input.trim());
        });
      });

      switch (choix) {
        case '1':
          await afficherInformationsUtilisateur();
          break;
        case '2':
          await modifierInformationsUtilisateur();
          break;
        case '3':
          await menuChallenges();
          break;
        case '4':
          await deconnexion();
          await loginMenu();
          return;
        case '5':
          console.log('Au revoir !');
          readline.close();
          process.exit(0);
        default:
          console.log('Choix invalide. Veuillez réessayer.');
      }
    }
  }

  async function afficherInformationsUtilisateur() {
    const url = `http://${host}:${port}/api/user-info`;
    try {
      const response = await axios.get(url, {
        headers: {
          'Cookie': `${socengai_username_cookie_value}; ${socengai_auth_cookie_value}`
        }
      });
      console.log('Informations de l\'utilisateur :', response.data);
    } catch (error) {
      console.error(`Erreur lors de la récupération des informations de l'utilisateur : ${error.message}`);
    }
  }

  // Supprimer les cookies et le fichier cookies.txt
  async function delete_cookies() {
    socengai_username_cookie_value = '';
    socengai_auth_cookie_value = '';
    const fs = require('fs');

    if (fs.existsSync('cookies.txt')) {
      try {
        fs.unlinkSync('cookies.txt');
      } catch (error) {
        console.error(`Erreur lors de la suppression des cookies : ${error.message}`);
      }
    }
  }

  // Demande au serveur de vérifier la validité des cookies (route /check-cookies)
  async function ask_for_cookies_confirmation() {
    const url = `http://${host}:${port}/api/check-cookies`;
    console.log(`Cookies à vérifier : ${socengai_username_cookie_value}, ${socengai_auth_cookie_value}`);
    try {
      const response = await axios.get(url, {
        headers: {
          'Cookie': `${socengai_username_cookie_value}; ${socengai_auth_cookie_value}`
        }
      });

      if (response.status === 200) {
        console.log('Cookies valides.');
        return true;
      }
      else {
        console.log('Cookies invalides ou session expirée.');
        return false;
      }
    } catch (error) {
      console.log(`Session indisponible. Veuillez vous reconnecter.`);
      console.error(`Erreur lors de la vérification des cookies : ${error.message}`);
      await delete_cookies();
      await loginMenu();
    }
  }

  async function check_cookies() {
    const fs = require('fs');
    if (fs.existsSync('cookies.txt')) {
      const cookies = await readCookies();
      if (cookies.length > 0) {
        socengai_username_cookie_value = cookies.find(c => c.startsWith('socengai-username='));
        socengai_auth_cookie_value = cookies.find(c => c.startsWith('socengai-auth='));
        if (socengai_username_cookie_value && socengai_auth_cookie_value) {
          return socengai_username_cookie_value.split('=')[1];
        }
      }
    }
    return ''
  }
    

  async function creerCompte() {
    const username = await new Promise((resolve) => {
      readline.question('Entrez votre nom d\'utilisateur : ', resolve);
    });
    
    const email = await new Promise((resolve) => {
      readline.question('Entrez votre email : ', resolve);
    });
    
    const password = await new Promise((resolve) => {
      readline.question('Entrez votre mot de passe : ', resolve);
    });

    try {
      const url = `http://${host}:${port}/create-user`;
      const response = await axios.post(url, 
        JSON.stringify({
          username: username,
          email: email,
          password: password
        }), 
        {
          headers: {
            'Content-Type': 'application/json'
          }
        }
      );

      if (response.status === 201) {
        console.log('Compte créé avec succès !');
        save_cookies(response.headers['set-cookie']);
      } else {
        console.log('Erreur lors de la création du compte :', response.data);
      }
    } catch (error) {
      if (error.response) {
        console.log(error.response.data);
      } else {
        console.log(`Erreur de connexion : ${error.message}`);
      }
    }
    await loginMenu();
  }

  async function readCookies() {
    const fs = require('fs').promises;
    try {
      const data = await fs.readFile('cookies.txt', 'utf8');
      const cookies = data.split('\n').map(line => line.trim()).filter(line => line);
      return cookies;
    } catch (error) {
      console.error(`Erreur lors de la lecture des cookies : ${error.message}`);
      return [];
    }
  }

  // Programme principal
  async function run() {
    await initialise();
    await loginMenu();
    readline.close();
  }

  // Définition de la commande principale
  const program = new Command();
  program
    .name('mon-client-cli')
    .description('Client CLI pour interagir avec une API distante')
    .version('1.0.0');

  program.parse(process.argv);

  if (process.argv.slice(2).length === 0) {
    run();
  }
}

main();
