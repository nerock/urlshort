package router

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nerock/urlshort/server"
	"github.com/nerock/urlshort/url"
)

type URLService interface {
	CreateURL(context.Context, string) (string, error)
	GetURL(context.Context, string) (string, error)
	DeleteURL(context.Context, string) error
	IncrementRedirectionCount(context.Context, string) error
	GetRedirectionCount(context.Context, string) (int, error)
}

type URLRequest struct {
	URL string
}

type URLResponse struct {
	URL      string
	ShortURL string
}

type URLCountResponse struct {
	ID    string
	Count int
}

type URLRouter struct {
	urlSvc URLService
}

func NewURLRouter(urlSvc URLService) URLRouter {
	return URLRouter{urlSvc: urlSvc}
}

func (ur URLRouter) Routes(r *chi.Mux) {
	r.Get("/{id}", ur.redirectTo)
	r.Route("/api/url", func(r chi.Router) {
		r.Post("/", ur.createURL)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", ur.getURL)
			r.Delete("/", ur.deleteURL)
			r.Get("/count", ur.getCount)
		})
	})
}

func (ur URLRouter) redirectTo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		server.RenderError(w, errors.New("could not read id"), http.StatusBadRequest)
	}

	longURL, err := ur.urlSvc.GetURL(r.Context(), id)
	switch {
	case errors.Is(err, url.ErrNotFound):
		server.RenderError(w, err, http.StatusNotFound)
		return
	case err != nil:
		server.RenderError(w, err, http.StatusInternalServerError)
		return
	}

	if err := ur.urlSvc.IncrementRedirectionCount(r.Context(), id); err != nil {
		log.Println(err)
	}

	http.Redirect(w, r, longURL, http.StatusTemporaryRedirect)
}

func (ur URLRouter) createURL(w http.ResponseWriter, r *http.Request) {
	var req URLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.RenderError(w, err, http.StatusBadRequest)
	}

	shortURL, err := ur.urlSvc.CreateURL(r.Context(), req.URL)
	switch {
	case errors.Is(err, url.ErrInvalidURL):
		server.RenderError(w, err, http.StatusBadRequest)
		return
	case err != nil:
		server.RenderError(w, err, http.StatusInternalServerError)
		return
	}

	server.RenderSuccess(w, URLResponse{req.URL, shortURL}, http.StatusCreated)
}

func (ur URLRouter) getURL(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		server.RenderError(w, errors.New("could not read id"), http.StatusBadRequest)
	}

	longURL, err := ur.urlSvc.GetURL(r.Context(), id)
	switch {
	case errors.Is(err, url.ErrNotFound):
		server.RenderError(w, err, http.StatusNotFound)
		return
	case err != nil:
		server.RenderError(w, err, http.StatusInternalServerError)
		return
	}

	server.RenderSuccess(w, URLResponse{longURL, id}, http.StatusCreated)
}

func (ur URLRouter) deleteURL(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		server.RenderError(w, errors.New("could not read id"), http.StatusBadRequest)
	}

	err := ur.urlSvc.DeleteURL(r.Context(), id)
	switch {
	case errors.Is(err, url.ErrNotFound):
		server.RenderError(w, err, http.StatusNotFound)
		return
	case err != nil:
		server.RenderError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (ur URLRouter) getCount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		server.RenderError(w, errors.New("could not read id"), http.StatusBadRequest)
	}

	count, err := ur.urlSvc.GetRedirectionCount(r.Context(), id)
	switch {
	case errors.Is(err, url.ErrNotFound):
		server.RenderError(w, err, http.StatusNotFound)
		return
	case err != nil:
		server.RenderError(w, err, http.StatusInternalServerError)
		return
	}

	server.RenderSuccess(w, URLCountResponse{id, count}, http.StatusCreated)
}
