package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"soceng-ai/internals/utils"
	"soceng-ai/internals/utils/prompts"
)

var (
	env_params = []string{}
	wsl_params = []string{"GOOS=windows", "GOARCH=amd64"}
)

func main() {
	const (
		main_input    = "./cmd/main"
		main_output   = "./bin/main"
		server_input  = "./cmd/server"
		server_output = "./bin/socengai-server"
	)
	var ext string

	// Si on est sous WSL, on ajoute les variables d'environnement pour la compilation
	if utils.We_are_on_WSL() {
		env_params = wsl_params
	}
	if runtime.GOOS == "windows" || utils.We_are_on_WSL() {
		ext = ".exe"
	}

	compile(main_output, ext, main_input)
	compile(server_output, ext, server_input)

	advise_considering_env()
}

func compile(output string, ext string, input string) {
	cmd := exec.Command("go", "build", "-o", output+ext, input)

	// Ajout des variables d'environnement pour la compilation
	cmd.Env = append(os.Environ(), env_params...)

	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(prompts.Error+"Erreur lors de la compilation de", input, ":", err)
		fmt.Println("Sortie :\n", string(outputBytes))
	} else {
		fmt.Println(prompts.Success + " Compilation de " + input + " dans " + output + ext + " réussie.")
	}
}

func advise_considering_env() {
	if runtime.GOOS == "windows" || utils.We_are_on_WSL() {
		fmt.Println(prompts.Info + "Sachant que vous êtes sur Windows ou WSL, pensez à lancer l'exécutable en .exe.")
	}
}
