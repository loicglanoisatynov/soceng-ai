package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	colors "soceng-ai/internals/utils/colors"
	prompts "soceng-ai/internals/utils/prompts"
	"strings"
)

var (
	port  string
	host  string
	https bool
)

type ctxKey struct{}

func getField(r *http.Request, index int) string {
	fields := r.Context().Value(ctxKey{}).([]string)
	return fields[index]
}

func Serve(w http.ResponseWriter, r *http.Request) {
	var allow []string
	for _, route := range routes {
		matches := route.Get_route_regex().FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 {
			// if r.Method != route.Get_route_method() {
			// allow = append(allow, route.Get_route_method())
			// continue
			// }
			ctx := context.WithValue(r.Context(), ctxKey{}, matches[1:])
			route.Get_route_handler()(w, r.WithContext(ctx))
			return
		}
	}

	if len(allow) > 0 {
		w.Header().Set("Allow", strings.Join(allow, ", "))
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}

	http.NotFound(w, r)
}

// clearTerminal efface la console pour un affichage propre
func clearTerminal() {
	fmt.Print("\033[H\033[2J")
}

// StartServer initialise et démarre le serveur HTTP
func StartServer(args []string) {
	parseArgs(args)

	http.HandleFunc("/", Serve)

	if https {
		fmt.Println("HTTPS not yet implemented.")
		os.Exit(1)
		// http.ListenAndServeTLS(":"+port, "cert.pem", "key.pem", nil)
	}
	if runtime.GOOS != "windows" && port < "1024" {
		fmt.Println("Démarrez le serveur avec sudo pour utiliser un port inférieur à 1024.")
		os.Exit(0)
	}
	fmt.Println(prompts.Prompt + prompts.Success + "Serveur HTTP démarré sur " + colors.Cyan + host + colors.Reset + ":" + colors.Cyan + port + colors.Reset)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func parseArgs(strings []string) {
	for i, s := range strings {
		if s == "-p" && i+1 < len(strings) {
			port = strings[i+1]
		} else if s == "-h" && i+1 < len(strings) {
			host = strings[i+1]
		} else if s == "-s" {
			https = true
		}
	}
	if port == "" {
		port = "80"
	}
	if host == "" {
		host = "localhost"
	}
}
