package services

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	env "soceng-ai/internals/server/env"
	utils "soceng-ai/internals/utils"
	"soceng-ai/internals/utils/colors"
	"soceng-ai/internals/utils/prompts"
)

var (
	port, host  string
	https, skip bool
)

// Lance le serveur en arrière-plan. Fonction principale du fichier start.go.
func Start(args []string) {
	if check_if_running() {
		fmt.Println("Le serveur est déjà en cours d'exécution.")
		return
	} else {
		fmt.Println(prompts.Prompt + "Lancement du serveur...")
	}

	parseArgs(args)

	attr := &os.ProcAttr{
		Dir: ".",
		Env: os.Environ(),
		Files: []*os.File{
			os.Stdin,
			os.Stdout,
			os.Stderr,
		},
		// Sys: &syscall.SysProcAttr{
		// 	HideWindow: true,
		// 	CmdLine:    "server",
		// CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
		// },
	}

	if runtime.GOOS == "windows" || utils.We_are_on_WSL() {
		processus += winexe
		if runtime.GOOS == "windows" {
			env.BINPATH = strings.ReplaceAll(env.BINPATH, "/", "\\")
		}
	}

	if !skip {
		init_server()
	} else {
		host = "localhost"
		if https {
			port = "443"
		} else {
			port = "80"
		}
	}
	serv_args := []string{env.BINPATH + processus, "-host", host, "-port", port}

	// Lancer le processus en arrière-plan
	proc, err := os.StartProcess(env.BINPATH+processus, serv_args, attr)
	if err != nil {
		log.Fatalf("Erreur lors du démarrage du serveur : %v\n", err)
	}

	fmt.Println("Serveur démarré en arrière-plan (PID :", proc.Pid, ")")

	err = os.WriteFile(env.PID_PATH, []byte(strconv.Itoa(proc.Pid)), 0644)
	if err != nil {
		log.Fatalf("Erreur lors de l'écriture du fichier PID : %v\n", err)
	}
	fmt.Println("PID du serveur écrit dans", env.PID_PATH)

	runtimeArgs := []string{"-host", host, "-port", port}
	cmd := exec.Command(env.BINPATH+processus, runtimeArgs...)
	err = cmd.Start()
	if err != nil {
		log.Fatalln(err)
	}
	// proc.Release()
}

func parseArgs(args []string) {
	for i := 0; i < len(args); i++ {
		if args[i] == "-h" || args[i] == "--help" {
			fmt.Println("Usage: server [OPTIONS]\n\nOPTIONS:")
			fmt.Println("  -h, --help\t\tAffiche ce message d'aide.")
			fmt.Println("  -https\t\t\tUtilise le protocole HTTPS.")
			fmt.Println("  -s, --skip\t\tIgnore les questions de configuration.")
			fmt.Println("  --host, -H\t\tSpécifie le HOST sur lequel le serveur écoutera.")
			fmt.Println("  --port, -p\t\tSpécifie le PORT sur lequel le serveur écoutera.")
			os.Exit(0)
		}
		if args[i] == "-https" {
			https = true
		}
		if args[i] == "-s" || args[i] == "--skip" {
			skip = true
		}
		if args[i] == "--host" || args[i] == "-H" {
			host = args[i+1]
		}
		if args[i] == "--port" || args[i] == "-p" {
			port = args[i+1]
		}
	}
}

func init_server() {
	not_valid := true
	for not_valid {
		fmt.Println(prompts.Prompt + "Initializing server...")
		if port == "" {
			define_port()
		}
		if host == "" {
			define_host()
		}
		not_valid = confirm_user_input()
	}
}

func define_port() {
	for port == "" {
		print(prompts.Prompt + "Enter the " + colors.Cyan_ify("PORT") + " on which the server will listen : ")
		var input string
		_, _ = fmt.Scanln(&input)
		port = input
		if port == "" {
			port = default_port
		} else if !is_valid_port() {
			println(prompts.Prompt + "Invalid " + colors.Cyan_ify("PORT") + ". Please enter a valid " + colors.Cyan_ify("PORT") + ".")
			port = ""
		} else {
			break
		}
	}

}

func define_host() {
	for host == "" {
		print(prompts.Prompt + "Enter the " + colors.Cyan_ify("HOST") + " on which the server will listen : ")
		var input string
		_, _ = fmt.Scanln(&input)
		host = input
		if host == "" {
			host = "localhost"
		} else if !is_valid_host(host) {
			println(prompts.Prompt + "Invalid " + colors.Cyan_ify("HOST") + ". Please enter a valid " + colors.Cyan_ify("HOST") + ".")
			host = ""
		} else {
			break
		}
	}
}

func is_valid_port() bool {
	port_int, err := strconv.Atoi(port)
	return err == nil && port_int > 1024 && port_int < 65536
}

func is_valid_host(host string) bool {
	u, err := url.Parse(host)
	return (err == nil && u.Scheme != "" && u.Host != "") || host == "localhost"
}

func confirm_user_input() bool {
	println(prompts.Prompt + "Server will listen on " + colors.Cyan_ify(host) + ":" + colors.Cyan_ify(port))
	print(prompts.Prompt + "Is this information correct? [y/n] : ")
	var input string
	_, _ = fmt.Scanln(&input)
	return input != "y"
}
