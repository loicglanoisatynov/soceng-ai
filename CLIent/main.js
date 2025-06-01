const { Command } = require('commander');
const readline = require('readline').createInterface({
  input: process.stdin,
  output: process.stdout
});
const axios = require('axios');
const { read } = require('fs');

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
    while (true) {
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
      }
    }
  }


  async function seConnecter() {
    
    let cookies_exist, username_cookie = await check_cookies();
    if (cookies_exist) {
    let reponse = await new Promise((resolve) => {
      readline.question('Une session est déjà active, au nom de : ' + username_cookie + '. Voulez-vous utiliser cette session ? (O/N)', (input) => {
        resolve(input.trim().toUpperCase());
      });
    });
      if (reponse === 'O') {
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
        save_cookies(response.headers['set-cookie']);
        dashboard_menu();
      } else {
        console.log('Erreur lors de la connexion :', response.data);
      }
    } catch (error) {
      if (error.response && error.response.data) {
        if (error.response.data === "Invalid password\n") {
          console.log('Mot de passe invalide. Veuillez réessayer.');
        } else {
          console.log(error.response.data);
        }
      } else {
        console.log(`Erreur de connexion : ${error.message}`);
      }
    }
  }

  async function dashboard_menu() {
    while (true) {
      console.log('Menu du tableau de bord :');
      console.log('1. Afficher les informations de l\'utilisateur');
      console.log('2. Modifier les informations de l\'utilisateur');
      console.log('3. Menu challenges');
      console.log('4. Déconnexion');
      console.log('5. Quitter');
      let choix = await new Promise((resolve) => {
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
          return; // Retour au menu de connexion
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
    try {
      const response = await axios.get(url, {
        headers: {
          'Cookie': `${socengai_username_cookie_value}; ${socengai_auth_cookie_value}`
        }
      });

      return response.data.valid;
    } catch (error) {
      console.log(`Session indisponible. Veuillez vous reconnecter.`);
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
          return true, socengai_username_cookie_value.split('=')[1];
        }
      }
    }
    return false, '';
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
    const fs = require('fs');
    fs.writeFile('cookies.txt', socengai_username_cookie_value, (err) => {
      if (err) {
        console.error(`Erreur lors de l'enregistrement des cookies : ${err.message}`);
      } else {
        console.log('Cookies enregistrés avec succès dans cookies.txt');
      }
    });
    fs.appendFile('cookies.txt', `\n${socengai_auth_cookie_value}`, (err) => {
      if (err) {
        console.error(`Erreur lors de l'ajout du cookie d'authentification : ${err.message}`);
      } else {
        console.log('Cookie d\'authentification ajouté avec succès dans cookies.txt');
      }
    });
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
