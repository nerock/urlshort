package server

import (
	"context"
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
