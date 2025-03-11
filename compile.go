package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	services "soceng-ai/ctrl-cmd/main/services"
	env "soceng-ai/internals/server/env"
	"soceng-ai/internals/utils"
	"soceng-ai/internals/utils/prompts"
)

const (
	main_input   = "./ctrl-cmd/main"
	main_output  = "./bin/ctrl-cmd"
	server_input = "./internals"
)

var (
	server_output = env.BINPATH + env.PROCESS
	env_params    = []string{}
	wsl_params    = []string{"GOOS=windows", "GOARCH=amd64"}
)

func main() {
	cleanBin()
	killProcess()

	var ext string

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
	env := ""
	if runtime.GOOS == "windows" || utils.We_are_on_WSL() {
		if utils.We_are_on_WSL() {
			env = "WSL"
		} else {
			env = "Windows"
		}
		fmt.Println(prompts.Info + "Sachant que vous êtes sur " + env + ", pensez à lancer l'exécutable en .exe.")
	}
}

func cleanBin() {
	os.RemoveAll("./bin")
	os.Mkdir("./bin", 0755)
}

func killProcess() {
	if services.Get_process_id_from_process_name() != 0 {
		fmt.Println(prompts.Info + "Arrêt du serveur en cours...")
		process, err := os.FindProcess(services.Get_process_id_from_process_name())
		if err != nil {
			fmt.Println(prompts.Error+"Erreur lors de la recherche du processus :", err)
		}
		err = process.Kill()
		if err != nil {
			fmt.Println(prompts.Error+"Erreur lors de l'arrêt du processus :", err)
		}
		fmt.Println(prompts.Success + "Serveur arrêté.")
	}
}
