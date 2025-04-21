package handlers

import (
	"net/http"
)

func Get(h http.HandlerFunc) http.HandlerFunc {
	return allowMethod(h, "GET")
}

func allowMethod(h http.HandlerFunc, method string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if method != r.Method {
			w.Header().Set("Allow", method)
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h(w, r)
	}
}

func Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the home page!\n"))
}

func Helloworld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!"))
}
