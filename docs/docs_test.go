package docs_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/nerock/urlshort/docs"
)

func TestGetDocs(t *testing.T) {
	tests := map[string]struct {
		wantStatus int
		wantBody   []byte
	}{
		"success": {
			wantStatus: http.StatusOK,
			wantBody:   []byte(docs.SwaggerUI),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			srv := httptest.NewServer(getRouter())
			res, err := http.Get(srv.URL + path.Join("/api/docs"))
			if err != nil {
				t.Errorf("could not send request: %v", err)
				return
			}

			checkResponse(t, res, tt.wantStatus, tt.wantBody)
		})
	}
}

func TestGetSwaggerJSON(t *testing.T) {
	docs.SwaggerJSON = "testJSON"

	tests := map[string]struct {
		wantStatus int
		wantBody   []byte
	}{
		"success": {
			wantStatus: http.StatusOK,
			wantBody:   []byte("testJSON"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			srv := httptest.NewServer(getRouter())
			res, err := http.Get(srv.URL + path.Join("/api/docs/swagger.json"))
			if err != nil {
				t.Errorf("could not send request: %v", err)
				return
			}

			checkResponse(t, res, tt.wantStatus, tt.wantBody)
		})
	}
}

func getRouter() *chi.Mux {
	r := chi.NewRouter()
	urlRouter := docs.Router{}
	urlRouter.Routes(r)

	return r
}

func checkResponse(t *testing.T, res *http.Response, expectedStatus int, expectedBody []byte) {
	t.Helper()

	if res.StatusCode != expectedStatus {
		t.Errorf("wrong status code returned\nexpected=%d\ngot=%d", expectedStatus, res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("could not read body: %s", err)
		return
	}

	body = bytes.TrimSpace(body)
	if !bytes.Equal(body, expectedBody) {
		t.Errorf("wrong body returned\nexpected=%s\ngot=%s", expectedBody, body)
	}

	if err := res.Body.Close(); err != nil {
		t.Errorf("could not close response body: %s", err)
	}
}
