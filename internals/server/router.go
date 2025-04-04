package server

import (
	"fmt"
	"net/http"
	"regexp"
	handlers "soceng-ai/internals/server/handlers"
	"soceng-ai/internals/server/handlers/registering"
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
		newRoute("GET", "/read-user", registering.Read_user),
		newRoute("PUT", "/update-user", registering.Update_user),
		newRoute("DELETE", "/delete-user", registering.Delete_user),

		// newRoute("GET", "/contact", contact),
		// newRoute("GET", "/api/widgets", apiGetWidgets),
		// newRoute("POST", "/api/widgets", apiCreateWidget),
		// newRoute("POST", "/api/widgets/([^/]+)", apiUpdateWidget),
		// newRoute("POST", "/api/widgets/([^/]+)/parts", apiCreateWidgetPart),
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
