/* Contenu du fichier soceng-ai/internals/server/server.go */
package server

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	// "github.com/shirou/gopsutil/cpu"
	// "github.com/shirou/gopsutil/mem"

	routes "soceng-ai/internals/server/routes"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

type ServerStats struct {
	mu              sync.Mutex
	activeConn      int
	totalRequests   int
	lastConnections []string
}

var stats = ServerStats{}
var port string
var host string
var https bool

func StartServer(args []string) {
	parseArgs(args)

	// go monitorServer()

	http.HandleFunc("/", Serve)

	if https {
		print("HTTPS not yet implemented.")
		os.Exit(1)
		// http.ListenAndServeTLS(":"+port, "cert.pem", "key.pem", nil)
	} else {
		http.ListenAndServe(host+":"+port, nil)
	}
}

/*
	Reprend l approche de Split Switching du blog benhoyt.com accessible à l adresse suivante :

https://benhoyt.com/writings/go-routing/#split-switch
*/
func Serve(w http.ResponseWriter, r *http.Request) {
	p := strings.Split(r.URL.Path, "/")[1:]
	n := len(p)

	// Mettre à jour les statistiques
	stats.mu.Lock()
	stats.activeConn++
	stats.totalRequests++
	if len(stats.lastConnections) >= 5 {
		stats.lastConnections = stats.lastConnections[1:] // Supprime la plus ancienne entrée
	}
	stats.lastConnections = append(stats.lastConnections, r.RemoteAddr)
	stats.mu.Unlock()
	// Fin de la mise à jour des statistiques

	var h http.Handler
	switch {
	case (n == 1 && p[0] == "") || (n == 1 && p[0] == "home"):
		h = routes.Get(routes.Home)
	case n == 1 && p[0] == "helloworld":
		h = routes.Get(routes.Helloworld)
	// case n == 1 && p[0] == "contact":
	// 	h = get(contact)
	// case n == 2 && p[0] == "api" && p[1] == "widgets" && r.Method == "GET":
	// 	h = get(apiGetWidgets)
	// case n == 2 && p[0] == "api" && p[1] == "widgets":
	// 	h = post(apiCreateWidget)
	// case n == 3 && p[0] == "api" && p[1] == "widgets" && p[2] != "":
	// 	h = post(apiWidget{p[2]}.update)
	// case n == 4 && p[0] == "api" && p[1] == "widgets" && p[2] != "" && p[3] == "parts":
	// 	h = post(apiWidget{p[2]}.createPart)
	// case n == 6 && p[0] == "api" && p[1] == "widgets" && p[2] != "" && p[3] == "parts" && isId(p[4], &id) && p[5] == "update":
	// 	h = post(apiWidgetPart{p[2], id}.update)
	// case n == 6 && p[0] == "api" && p[1] == "widgets" && p[2] != "" && p[3] == "parts" && isId(p[4], &id) && p[5] == "delete":
	// 	h = post(apiWidgetPart{p[2], id}.delete)
	// case n == 1:
	// 	h = get(widget{p[0]}.widget)
	// case n == 2 && p[1] == "admin":
	// 	h = get(widget{p[0]}.admin)
	// case n == 2 && p[1] == "image":
	// 	h = post(widget{p[0]}.image)
	default:
		http.NotFound(w, r)
		return
	}
	h.ServeHTTP(w, r)

	stats.mu.Lock()
	stats.activeConn-- // Une fois la requête traitée, on diminue le compteur
	stats.mu.Unlock()
}

func monitorServer() {
	for {
		stats.mu.Lock()
		clearTerminal() // Nettoie la console pour un affichage dynamique

		// Récupérer l'utilisation CPU et RAM
		vMem, _ := mem.VirtualMemory()
		cpuPercent, _ := cpu.Percent(time.Second, false)

		// Afficher les statistiques
		fmt.Println("===== MONITORING DU SERVEUR =====")
		fmt.Printf("🔹 Connexions actives  : %d\n", stats.activeConn)
		fmt.Printf("🔹 Requêtes totales    : %d\n", stats.totalRequests)
		fmt.Printf("🔹 Dernières IPs       : %v\n", stats.lastConnections)
		fmt.Printf("🔹 CPU Usage          : %.2f%%\n", cpuPercent[0])
		fmt.Printf("🔹 RAM Usage          : %.2f%%\n", vMem.UsedPercent)
		fmt.Printf("🔹 Go Routines        : %d\n", runtime.NumGoroutine())
		fmt.Println("================================")

		stats.mu.Unlock()

		time.Sleep(10 * time.Second) // Rafraîchissement toutes les 2 secondes
	}
}

// Fonction pour nettoyer la console (pour un affichage propre)
func clearTerminal() {
	fmt.Print("\033[H\033[2J")
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
