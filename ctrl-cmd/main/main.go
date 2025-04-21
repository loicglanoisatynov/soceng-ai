package main

import (
	"fmt"
	"os"
	services "soceng-ai/ctrl-cmd/main/services"
)

/* Point d'entrée du programme. Sert à démarrer le serveur. Server.exe (cmd/server.main.go) ne devrait pas être démarré manuellement. */
func main() {

	if len(os.Args) < 2 {
		printhelp()
	}

	switch os.Args[1] {
	case "start":
		services.Start(os.Args)
	case "stop":
		services.Stop()
	// case "status":
	// 	Status()
	// case "dashboard":
	// 	Dashboard()
	default:
		printhelp()
	}
}

/* Affiche l'aide du serveur. */
func printhelp() {
	fmt.Println("Usage : serv <start|stop|status|dashboard>")
	os.Exit(1)
}
