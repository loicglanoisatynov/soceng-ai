package main

import (
	"fmt"
	"net/http"
	"os"
	database "soceng-ai/database"
	_ "soceng-ai/internals/docs" // Importez le package docs généré par Swaggo
	"soceng-ai/internals/server"
	"soceng-ai/internals/utils/prompts"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"
)

var (
	dev_mode   bool
	port, host string
)

func main() {
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	if len(os.Args) < 2 {
		os.Args = append(os.Args, "-h")
		os.Args = append(os.Args, "127.0.0.1")
		os.Args = append(os.Args, "-p")
		os.Args = append(os.Args, "80")
	}

	parseArgs(os.Args[1:])

	database.Init_DB()
	server.StartServer(os.Args[1:])
}

func Get_prompt() string {
	return "\033[36m[ " + time.Now().Format("2006-01-02 15:04:05.000000") + " server ]\033[0m "
}

func parseArgs(args []string) {
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-h":
			host = args[i+1]
		case "-p":
			port = args[i+1]
		case "-d", "--dev", "--devmode=on":
			dev_mode = true
			fmt.Println(prompts.Debug + "Mode développeur activé.")
		}
	}
}
