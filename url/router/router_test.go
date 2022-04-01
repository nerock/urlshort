package router_test

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/nerock/urlshort/url"
	"github.com/nerock/urlshort/url/router"
)

var errSvc = errors.New("service error")

type testService struct {
	id    string
	url   string
	count int
	err   error
}

func (t testService) CreateURL(ctx context.Context, s string) (string, error) {
	return t.id, t.err
}

func (t testService) GetURL(ctx context.Context, s string) (string, string, error) {
	return t.url, t.id, t.err
}

func (t testService) DeleteURL(ctx context.Context, s string) error {
	return t.err
}

func (t *testService) IncrementRedirectionCount(ctx context.Context, s string) error {
	t.count++
	return t.err
}

func (t testService) GetRedirectionCount(ctx context.Context, s string) (int, error) {
	return t.count, t.err
}

func TestRedirect(t *testing.T) {
	tests := map[string]struct {
		testSvc testService

		wantStatus int
		wantBody   []byte
	}{
		"id not found": {
			testSvc: testService{
				err: url.ErrNotFound,
			},
			wantStatus: http.StatusNotFound,
			wantBody:   []byte(`{"Code":"Not Found","Message":"` + url.ErrNotFound.Error() + `"}`),
		},
		"svc error": {
			testSvc: testService{
				err: errSvc,
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   []byte(`{"Code":"Internal Server Error","Message":"` + errSvc.Error() + `"}`),
		},
		"success": {
			testSvc: testService{
				count: 10,
				url:   "https://www.google.es",
			},
			wantStatus: http.StatusOK,
			wantBody:   []byte(`{"ID":"ID","Count":10}`),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			initialCount := tt.testSvc.count

			srv := httptest.NewServer(getRouter(&tt.testSvc))
			res, err := http.Get(srv.URL + path.Join("/ID"))
			if err != nil {
				t.Errorf("could not send request: %v", err)
				return
			}

			if res.StatusCode != http.StatusOK {
				if initialCount != tt.testSvc.count {
					t.Errorf("redirection count should not change\nexpected=%d\ngot=%d", initialCount, tt.testSvc.count)
				}

				checkResponse(t, res, tt.wantStatus, tt.wantBody)
			} else { // If status code was 200 it redirected correctly
				if initialCount+1 != tt.testSvc.count {
					t.Errorf("redirection count should have incremented by one\nexpected=%d\ngot=%d", initialCount+1, tt.testSvc.count)
				}
			}
		})
	}
}

func TestCreateURL(t *testing.T) {
	tests := map[string]struct {
		testSvc     testService
		requestBody []byte

		wantStatus int
		wantBody   []byte
	}{
		"invalid url": {
			requestBody: []byte(`{"URL":"url"}`),
			testSvc: testService{
				err: url.ErrInvalidURL,
				url: "url",
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   []byte(`{"Code":"Bad Request","Message":"` + url.ErrInvalidURL.Error() + `"}`),
		},
		"svc error": {
			requestBody: []byte(`{"URL":"url"}`),
			testSvc: testService{
				err: errSvc,
				url: "url",
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   []byte(`{"Code":"Internal Server Error","Message":"` + errSvc.Error() + `"}`),
		},
		"success": {
			requestBody: []byte(`{"URL":"url"}`),
			testSvc: testService{
				id:  "id",
				url: "url",
			},
			wantStatus: http.StatusCreated,
			wantBody:   []byte(`{"URL":"url","ShortURL":"id"}`),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			srv := httptest.NewServer(getRouter(&tt.testSvc))
			res, err := http.Post(srv.URL+path.Join("/api/url"), "application/json",
				bytes.NewReader(tt.requestBody))
			if err != nil {
				t.Errorf("could not send request: %v", err)
				return
			}

			checkResponse(t, res, tt.wantStatus, tt.wantBody)
		})
	}
}

func TestGetURL(t *testing.T) {
	tests := map[string]struct {
		testSvc testService

		wantStatus int
		wantBody   []byte
	}{
		"id not found": {
			testSvc: testService{
				err: url.ErrNotFound,
			},
			wantStatus: http.StatusNotFound,
			wantBody:   []byte(`{"Code":"Not Found","Message":"` + url.ErrNotFound.Error() + `"}`),
		},
		"svc error": {
			testSvc: testService{
				err: errSvc,
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   []byte(`{"Code":"Internal Server Error","Message":"` + errSvc.Error() + `"}`),
		},
		"success": {
			testSvc: testService{
				url: "url",
				id:  "ID",
			},
			wantStatus: http.StatusOK,
			wantBody:   []byte(`{"URL":"url","ShortURL":"ID"}`),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			srv := httptest.NewServer(getRouter(&tt.testSvc))
			res, err := http.Get(srv.URL + path.Join("/api/url/ID"))
			if err != nil {
				t.Errorf("could not send request: %v", err)
				return
			}

			checkResponse(t, res, tt.wantStatus, tt.wantBody)
		})
	}
}

func TestDeleteURL(t *testing.T) {
	tests := map[string]struct {
		testSvc testService

		wantStatus int
		wantBody   []byte
	}{
		"id not found": {
			testSvc: testService{
				err: url.ErrNotFound,
			},
			wantStatus: http.StatusNotFound,
			wantBody:   []byte(`{"Code":"Not Found","Message":"` + url.ErrNotFound.Error() + `"}`),
		},
		"svc error": {
			testSvc: testService{
				err: errSvc,
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   []byte(`{"Code":"Internal Server Error","Message":"` + errSvc.Error() + `"}`),
		},
		"success": {
			testSvc: testService{
				url: "url",
			},
			wantStatus: http.StatusNoContent,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			srv := httptest.NewServer(getRouter(&tt.testSvc))
			req, err := http.NewRequest(http.MethodDelete, srv.URL+path.Join("/api/url/ID"), nil)
			if err != nil {
				t.Errorf("could not create request: %s", err)
				return
			}
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Errorf("could not send request: %v", err)
				return
			}

			checkResponse(t, res, tt.wantStatus, tt.wantBody)
		})
	}
}

func TestGetCount(t *testing.T) {
	tests := map[string]struct {
		testSvc testService

		wantStatus int
		wantBody   []byte
	}{
		"id not found": {
			testSvc: testService{
				err: url.ErrNotFound,
			},
			wantStatus: http.StatusNotFound,
			wantBody:   []byte(`{"Code":"Not Found","Message":"` + url.ErrNotFound.Error() + `"}`),
		},
		"svc error": {
			testSvc: testService{
				err: errSvc,
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   []byte(`{"Code":"Internal Server Error","Message":"` + errSvc.Error() + `"}`),
		},
		"success": {
			testSvc: testService{
				count: 10,
			},
			wantStatus: http.StatusOK,
			wantBody:   []byte(`{"ID":"ID","Count":10}`),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			srv := httptest.NewServer(getRouter(&tt.testSvc))
			res, err := http.Get(srv.URL + path.Join("/api/url/ID/count"))
			if err != nil {
				t.Errorf("could not send request: %v", err)
				return
			}

			checkResponse(t, res, tt.wantStatus, tt.wantBody)
		})
	}
}

func getRouter(svc router.URLService) *chi.Mux {
	r := chi.NewRouter()
	urlRouter := router.NewURLRouter(svc)
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
