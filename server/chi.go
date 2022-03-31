package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type router interface {
	Routes(*chi.Mux)
}

type HTTPServer struct {
	router *chi.Mux
	srv    *http.Server
}

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

func (s HTTPServer) Run(routers ...router) error {
	for _, r := range routers {
		r.Routes(s.router)
	}

	return s.srv.ListenAndServe()
}

func (s HTTPServer) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func RenderSuccess(w http.ResponseWriter, res any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "could not encode response", http.StatusInternalServerError)
	}
}

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
