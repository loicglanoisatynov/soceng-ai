Projet SocEng-AI - README

# Sommaire
- [Sommaire](#sommaire)
- [Introduction](#introduction)
- [Installation](#installation)
  - [Etape 1 : Cloner le dépôt](#etape-1--cloner-le-dépôt)
  - [Etape 2 : Installer les dépendances](#etape-2--installer-les-dépendances)
- [Lancement](#lancement)
- [Contributeurs](#contributeurs)
- [Licence](#licence)
- [Contact](#contact)

# Introduction

Ce projet consiste en la création d'une plateforme de CTFs basés sur le principe de l'ingénierie sociale, où les participants doivent interagir avec des personnages IA pour mener à bien des missions, où on obtient des documents et des contacts avec d'autres personnages IA pour progresser dans le jeu et avancer vers le flag capital, qui conclus le challenge. Chaque personnage IA a son propre lore ainsi qu'un niveau de suspicion qui varie avec les interactions de l'utilisateur. Le challenge est perdu si le niveau de suspicion atteint une valeur de 10. Ce projet est développé dans le cadre du projet dev "fil rouge" du cursus de B2 Info de Lyon Ynov Campus.

# Installation

## Etape 1 : Cloner le dépôt

Pour installer le projet, il est nécessaire de cloner le dépôt Git et d'installer les dépendances requises. Voici les étapes à suivre :

1. Installer l'utilitaire `git` si ce n'est pas déjà fait. Suivez la procédure d'installation pour votre système d'exploitation :
   - Pour **Linux** : Utilisez votre gestionnaire de paquets (par exemple, `apt` pour Debian/Ubuntu, `yum` pour CentOS, etc.).
   - Pour **macOS** : Installez via Homebrew avec la commande `brew install git`.
   - Pour **Windows** : Téléchargez et installez Git depuis [git-scm.com](https://git-scm.com/download/win).

    Pour plus d'informations, consultez la [documentation officielle de Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git).

2. Ouvrir un terminal ou une invite de commande :
   - Sur **Linux** ou **macOS**, ouvrez l'application Terminal.
   - Sur **Windows**, ouvrez l'invite de commande ou PowerShell.
  
3. Naviguez vers le répertoire où vous souhaitez cloner le dépôt. Par exemple, pour aller dans le dossier `Documents` :
   - Sur **Linux** ou **macOS** :
    ```bash
    cd /home/username/Documents
    ```
   - Sur **Windows** :
    ```cmd
    cd C:\Users\username\Documents
    ```

4. Cloner le dépôt Git du projet SocEng-AI :
    ```bash
    git clone https://github.com/loicglanoisatynov/soceng-ai.git
    ```

5. Accéder au répertoire du projet cloné :
    ```bash
    cd soceng-ai
    ```

## Etape 2 : Installer les dépendances

Le projet nécessite les dépendances suivantes pour fonctionner correctement :
- `go`: Le langage de programmation utilisé pour développer le projet.
- `node`: Utilisé pour le client en ligne de commandes.
- `typescript`: Pour le développement du frontend.
- `angular`: Framework pour le développement du frontend.
- `tailwindcss`: Utilisé pour le style du frontend.
- `ng`: Outil de ligne de commande pour Angular.

# Lancement

Pour lancer le serveur du projet, suivez ces étapes :
1. Assurez-vous d'avoir installé les dépendances nécessaires comme indiqué dans la section précédente.
2. Ouvrez un terminal ou une invite de commande.
3. Naviguez vers le répertoire du projet cloné :
   ```bash
   cd /chemin/vers/soceng-ai
   ```
4. Lancez le serveur avec la commande suivante :
   ```bash
    go run internals/main.go
    ```

# Utilisation

En date du 4 juin 2025, seulement l'utilitaire en ligne de commande est fonctionnel. Pour l'utiliser, suivez ces étapes :
1. Ouvrez un terminal ou une invite de commande.
2. Naviguez vers le répertoire du projet cloné :
   ```bash
   node soceng-ai/CLIent/main.js
   ```
3. Suivez les instructions affichées pour interagir avec les personnages IA et progresser dans le jeu.

# Contributeurs

Ce projet a été développé par QUAGLIERI Lisa, KOUYATE Chouaib, DELPREE Corentin et Loïc GLANOIS, étudiants en B2 Info à Ynov Campus Lyon. Pour toute question ou contribution, n'hésitez pas à nous contacter via le dépôt GitHub.

# Licence

Ce projet est sous licence GPL-3.0. Vous pouvez librement utiliser, modifier et distribuer le code, à condition de respecter les termes de la licence. Pour plus de détails, consultez le fichier `LICENSE` dans le dépôt.

# Contact

Pour toute question ou suggestion concernant ce projet, vous pouvez nous contacter via le dépôt GitHub ou par email à l'adresse suivante : loic.glanois@ynov.com