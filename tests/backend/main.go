package main

import (
	"os"
	"soceng-ai/internals/utils/colors"
	"soceng-ai/internals/utils/prompts"
	"soceng-ai/tests/backend/auth"
	basics "soceng-ai/tests/backend/basics"
	"time"
)

var (
	port = "80"        // Port par défaut pour le serveur
	host = "localhost" // Host par défaut pour le serveur
)

// Fonction principale du processus de test du serveur backend.
func main() {
	parse_args()
	basics.Tests(host, port)
	auth.Tests(host, port)
	prompts.Prompts_tests(time.Now(), prompts.Success+"Tests du serveur backend terminés avec succès !")
}

// Traitement des arguments de la ligne de commande pour configurer le serveur backend.
func parse_args() {
	for i := 0; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-p", "--port":
			if i+1 < len(os.Args) {
				port = os.Args[i+1]
			}
		case "-h", "--host":
			if i+1 < len(os.Args) {
				host = os.Args[i+1]
			}
		}
	}

	prompts.Prompts_tests(time.Now(), prompts.Success+"Test du serveur backend sur l'hôte "+colors.Cyan_ify(host)+" et le port "+colors.Cyan_ify(port)+".")
}
