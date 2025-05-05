package server

import (
	"fmt"
	"net/http"
	"regexp"
	handlers "soceng-ai/internals/server/handlers"

	api "soceng-ai/internals/server/handlers/api"
	authentification "soceng-ai/internals/server/handlers/authentification"
	profiles_handling "soceng-ai/internals/server/handlers/profiles_handling"
	registering "soceng-ai/internals/server/handlers/registering"
)

var routes []Route

func routes_index_handler(w http.ResponseWriter, r *http.Request) {
	for _, route := range routes {
		fmt.Fprintf(w, "%s %s\n", route.Get_route_method(), route.Get_route_regex())
	}
}

func init() {
	routes = []Route{
		newRoute("GET", "/routes", routes_index_handler),
		newRoute("GET", "/", handlers.Home),
		newRoute("POST", "/create-user", registering.Register_user),
		newRoute("POST", "/check-register", registering.Check_register),
		newRoute("DELETE", "/delete-user", registering.Delete_user),
		newRoute("POST", "/login", authentification.Login),
		newRoute("DELETE", "/logout", authentification.Logout),
		newRoute("PUT", "/edit-profile", profiles_handling.Edit_profile),
		newRoute("PUT", "/edit-user", profiles_handling.Edit_user),

		// newRoute("GET", "/api/get-challenges", handlers.Get_challenges), // Récupère la liste des défis (notamment pour le front-end)
		newRoute("POST", "/api/challenge", api.Challenge_handler),
		newRoute("PUT", "/api/challenge", api.Challenge_handler),
		// newRoute("GET", "/api/get-challenge", handlers.Get_challenge),
		// newRoute("PUT", "/api/edit-challenge", handlers.Edit_challenge),
		// newRoute("DELETE", "/api/delete-challenge", handlers.Delete_challenge),

		// newRoute("GET", "/contact", contact),
		// newRoute("GET", "/([^/]+)/admin", widgetAdmin),
		// newRoute("POST", "/([^/]+)/image", widgetImage),
	}
}

// newRoute("GET", "/([^/]+)/admin", widgetAdmin),
// newRoute("POST", "/([^/]+)/image", widgetImage),

func Set_routes(new_routes []Route) {
	routes = new_routes
}

func Get_routes() []Route {
	return routes
}

func newRoute(method, pattern string, handler http.HandlerFunc) Route {
	return Route{method, regexp.MustCompile("^" + pattern + "$"), handler}
}

type Route struct {
	method  string
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

func (r Route) Get_route_method() string {
	return r.method
}

func (r Route) Get_route_regex() *regexp.Regexp {
	return r.regex
}

func (r Route) Get_route_handler() http.HandlerFunc {
	return r.handler
}
