package handlers

import (
    "net/http"
)

// WithCORS est un middleware qui gère les en-têtes CORS et les pré-vol (OPTIONS).
// Il prend un ServeMux (ou n'importe quel http.Handler) et renvoie un Handler
// qui injecte les bons headers avant d’appeler le handler d’origine.
func WithCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        origin := r.Header.Get("Origin")
        if origin != "" {
            w.Header().Set("Access-Control-Allow-Origin", origin)
            w.Header().Set("Vary", "Origin")
            w.Header().Set("Access-Control-Allow-Credentials", "true")
            w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        }

        if r.Method == http.MethodOptions {
            // Pré-vol : on répond tout de suite
            w.WriteHeader(http.StatusOK)
            return
        }

        // Sinon, on continue vers le vrai handler
        next.ServeHTTP(w, r)
    })
}
