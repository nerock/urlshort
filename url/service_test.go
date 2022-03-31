package url_test

import (
	"context"
	"errors"
	"testing"

	"github.com/nerock/urlshort/url"
)

var (
	errGenerator = errors.New("generator error")
	errStore     = errors.New("store error")
)

type testGenerator struct {
	id  string
	err error
}

func (t testGenerator) Generate() (string, error) {
	return t.id, t.err
}

type testStore struct {
	url string
	err error
}

func (t testStore) AddURL(ctx context.Context, short, long string) error {
	return t.err
}

func (t testStore) GetURL(ctx context.Context, short string) (string, error) {
	return t.url, t.err
}

func (t testStore) DeleteURL(ctx context.Context, short string) error {
	return t.err
}

func TestCreate(t *testing.T) {
	validURL := "https://www.google.es"
	invalidURL := "invalidURL"

	domain := "localhost:8080/"

	validID := "ID"
	invalidID := ""

	tests := map[string]struct {
		store     testStore
		generator testGenerator

		url string

		id  string
		err error
	}{
		"generator error": {
			store: testStore{
				url: validURL,
			},
			generator: testGenerator{
				id:  invalidID,
				err: errGenerator,
			},
			url: validURL,
			err: errGenerator,
		},
		"store error": {
			store: testStore{
				err: errStore,
			},
			generator: testGenerator{
				id: validID,
			},
			url: validURL,
			err: errStore,
		},
		"invalid URL": {
			store: testStore{
				url: validURL,
			},
			generator: testGenerator{
				id: validID,
			},
			url: invalidURL,
			err: url.ErrInvalidURL,
		},
		"success": {
			store: testStore{
				url: validURL,
			},
			generator: testGenerator{
				id: validID,
			},
			url: validURL,
			id:  domain + validID,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			svc := url.NewService(domain, tt.generator, tt.store)
			id, err := svc.CreateURL(context.Background(), tt.url)

			if !errors.Is(err, tt.err) {
				t.Errorf("wrong error returned\nexpected=%s\ngot=%s", tt.err, err)
			}

			if id != tt.id {
				t.Errorf("wrong id returned\nexpected=%s\ngot=%s", tt.id, id)
			}
		})
	}
}

func TestGet(t *testing.T) {
	long := "https://www.google.es"
	short := "ID"

	tests := map[string]struct {
		store testStore

		short string

		url string
		err error
	}{
		"store error": {
			store: testStore{
				err: errStore,
			},
			short: short,
			url:   "",
			err:   errStore,
		},
		"success": {
			store: testStore{
				url: long,
			},
			short: short,
			url:   long,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			svc := url.NewService("", nil, tt.store)
			long, err := svc.GetURL(context.Background(), tt.url)

			if !errors.Is(err, tt.err) {
				t.Errorf("wrong error returned\nexpected=%s\ngot=%s", tt.err, err)
			}

			if long != tt.url {
				t.Errorf("wrong id returned\nexpected=%s\ngot=%s", tt.url, long)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	tests := map[string]struct {
		store testStore

		err error
	}{
		"store error": {
			store: testStore{
				err: errStore,
			},
			err: errStore,
		},
		"success": {
			err: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			svc := url.NewService("", nil, tt.store)
			err := svc.DeleteURL(context.Background(), "")

			if !errors.Is(err, tt.err) {
				t.Errorf("wrong error returned\nexpected=%s\ngot=%s", tt.err, err)
			}
		})
	}
}
