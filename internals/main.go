package handlers

import (
    "fmt"
    "net/http"
    "os"
    database "soceng-ai/database"
    server "soceng-ai/internals/server"
    handlers "soceng-ai/internals/server/handlers"
    "time"
)

var (
    dev_mode   bool
    port, host string
)

func main() {
    if len(os.Args) < 2 {
        os.Args = append(os.Args, "-h", "127.0.0.1", "-p", "80")
    }
    parseArgs(os.Args)

    database.Init_DB()
    server.StartServer(os.Args)

    addr := fmt.Sprintf("%s:%s", host, port)
    fmt.Printf("🚀 Listening on %s (CORS enabled)\n", addr)
    // On wrappe DefaultServeMux (qui contient toutes tes routes) avec le middleware CORS
    if err := http.ListenAndServe(addr, handlers.WithCORS(http.DefaultServeMux)); err != nil {
        panic("Serveur arrêté : " + err.Error())
    }
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
            fmt.Println("[DEBUG] Mode développeur activé.")
        }
    }
}
