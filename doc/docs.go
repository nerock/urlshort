package doc

import (
	_ "embed"
	"net/http"

	"github.com/go-chi/chi/v5"
)

//go:embed swagger.json
var swaggerJSON string

type Router struct{}

func (dr Router) Routes(r *chi.Mux) {
	r.Route("/docs", func(r chi.Router) {
		r.Get("/", dr.serveDocs)
		r.Get("/swagger.json", dr.serveJSON)
	})
}

func (Router) serveDocs(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(swaggerUI)); err != nil {
		http.Error(w, "could not render docs", http.StatusInternalServerError)
	}
}

func (Router) serveJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(swaggerJSON)); err != nil {
		http.Error(w, "could not render swagger JSON", http.StatusInternalServerError)
	}
}
