package server

import (
	"net/http"
	"strings"

	"socengai-server/server/routes"
)

// go env -w GOBIN=/mnt/c/Users/lglan/Laboratoire/Scolaire/B2Info-CampusYnov/Projet_Dev/SocEng.AI

// saphintosh@HP-Pavillon-Laptop-14-dv2xxx-Servitor:/mnt/c/Users/lglan/Laboratoire/Scolaire/B2Info-CampusYnov/Projet_Dev/SocEng.AI$ t
// ree
// .
// ├── README.md
// └── server
//     ├── go.mod
//     ├── routes
//     │   └── routes.go
//     └── soceng.ai-server.go

// 3 directories, 4 files

func Serve(w http.ResponseWriter, r *http.Request) {
	// Split path into slash-separated parts, for example, path "/foo/bar"
	// gives p==["foo", "bar"] and path "/" gives p==[""].
	p := strings.Split(r.URL.Path, "/")[1:]
	n := len(p)

	var h http.Handler
	// var id int
	switch {
	case n == 1 && p[0] == "":
		h = routes.Get(Home)
	case n == 1 && p[0] == "helloworld":
		h = routes.Get(helloworld)
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
}
