package main

import (
	"fmt"
	"os"
	"soceng-ai/internals/server"
	"soceng-ai/internals/utils/prompts"
	"time"
)

var (
	dev_mode   bool
	port, host string
)

func main() {

	if len(os.Args) < 2 {
		os.Args = append(os.Args, "-h")
		os.Args = append(os.Args, "localhost")
		os.Args = append(os.Args, "-p")
		os.Args = append(os.Args, "80")
	}

	parseArgs(os.Args)

	server.StartServer(os.Args)
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
