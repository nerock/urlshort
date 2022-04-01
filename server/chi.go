package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Router represents types that are able to add to the argument provided router
type Router interface {
	Routes(*chi.Mux)
}

// HTTPServer represents an HTTP Server
type HTTPServer struct {
	router *chi.Mux
	srv    *http.Server
}

// NewHTTPServer creates a new HTTPServer
func NewHTTPServer(port int) HTTPServer {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	return HTTPServer{
		router: r,
		srv: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: r,
		},
	}
}

// Run adds all the routes provided and then runs the HTTPServer
func (s HTTPServer) Run(routers ...Router) error {
	for _, r := range routers {
		r.Routes(s.router)
	}

	log.Println("Running HTTP server on:", s.srv.Addr)

	return s.srv.ListenAndServe()
}

// Shutdown gracefully shuts down the HTTPServer
func (s HTTPServer) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

// RenderSuccess renders a successful JSON response
func RenderSuccess(w http.ResponseWriter, res any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "could not encode response", http.StatusInternalServerError)
	}
}

// RenderError renders an error as JSON
func RenderError(w http.ResponseWriter, err error, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	res := struct {
		Code    string
		Message string
	}{
		Code:    http.StatusText(code),
		Message: err.Error(),
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "could not encode response", http.StatusInternalServerError)
	}
}
