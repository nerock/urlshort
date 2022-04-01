package docs

import (
	_ "embed"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// SwaggerJSON contains the api definition
//go:embed swagger.json
var SwaggerJSON string

// Router is the documentation router
type Router struct{}

// Routes adds documentation routes to the main router
func (dr Router) Routes(r *chi.Mux) {
	r.Route("/api/docs", func(r chi.Router) {
		r.Get("/", dr.serveDocs)
		r.Get("/swagger.json", dr.serveJSON)
	})
}

func (Router) serveDocs(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(SwaggerUI)); err != nil {
		http.Error(w, "could not render docs", http.StatusInternalServerError)
	}
}

func (Router) serveJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(SwaggerJSON)); err != nil {
		http.Error(w, "could not render swagger JSON", http.StatusInternalServerError)
	}
}
