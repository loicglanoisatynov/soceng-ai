package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

)

var (
	port  string
	host  string
	https bool
)

type ctxKey struct{}

// getField extrait un paramètre capturé par la regex dans le contexte de la requête
func getField(r *http.Request, index int) string {
	fields := r.Context().Value(ctxKey{}).([]string)
	return fields[index]
}

// Serve est le routeur principal (avec CORS géré en middleware)
func Serve(w http.ResponseWriter, r *http.Request) {
	// Le CORS est désormais géré en amont par handlers.WithCORS
	// =======================

	var allow []string
	for _, route := range routes {
		matches := route.Get_route_regex().FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 {
			if r.Method != route.Get_route_method() {
				allow = append(allow, route.Get_route_method())
				continue
			}

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

	// On enregistre désormais Serve sur DefaultServeMux
	http.HandleFunc("/", Serve)

	if https {
		fmt.Println("HTTPS not yet implemented.")
		os.Exit(1)
	}

	addr := host + ":" + port
	fmt.Println("Serveur HTTP démarré sur", addr)
	// Le ListenAndServe final est appelé dans main.go pour injecter le middleware CORS
}

// parseArgs lit les arguments pour le port, l'hôte et le mode HTTPS
func parseArgs(args []string) {
	for i, s := range args {
		switch s {
		case "-p":
			if i+1 < len(args) {
				port = args[i+1]
			}
		case "-h":
			if i+1 < len(args) {
				host = args[i+1]
			}
		case "-s":
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
