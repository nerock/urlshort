package url

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"path"
)

var (
	ErrInvalidURL = errors.New("invalid URL provided")
	ErrNotFound   = errors.New("URL not found")
)

type Generator interface {
	Generate() (string, error)
}

type Store interface {
	AddURL(ctx context.Context, short, long string) error
	GetURL(ctx context.Context, short string) (string, error)
	DeleteURL(ctx context.Context, short string) error
	IncrementRedirectionCount(ctx context.Context, short string) error
	GetRedirectionCount(ctx context.Context, short string) (int, error)
}

type Service struct {
	store     Store
	generator Generator

	domain string
}

func NewService(domain string, urlGenerator Generator, store Store) Service {
	return Service{
		domain:    domain,
		store:     store,
		generator: urlGenerator,
	}
}

func (s Service) CreateURL(ctx context.Context, long string) (string, error) {
	if _, err := url.ParseRequestURI(long); err != nil {
		return "", ErrInvalidURL
	}

	short, err := s.generator.Generate()
	if err != nil {
		return "", fmt.Errorf("could not generate URL: %w", err)
	}

	if err := s.store.AddURL(ctx, short, long); err != nil {
		return "", fmt.Errorf("could not save URL in database: %w", err)
	}

	return path.Join(s.domain, short), nil
}

func (s Service) GetURL(ctx context.Context, short string) (string, error) {
	long, err := s.store.GetURL(ctx, short)
	if err != nil {
		if err == ErrNotFound {
			return "", ErrNotFound
		}

		return "", fmt.Errorf("could not retrieve URL from database: %w", err)
	}

	return long, nil
}

func (s Service) DeleteURL(ctx context.Context, short string) error {
	if err := s.store.DeleteURL(ctx, short); err != nil {
		if err == ErrNotFound {
			return ErrNotFound
		}

		return fmt.Errorf("could not delete URL from database: %w", err)
	}

	return nil
}

func (s Service) IncrementRedirectionCount(ctx context.Context, short string) error {
	if err := s.store.IncrementRedirectionCount(ctx, short); err != nil {
		if err == ErrNotFound {
			return ErrNotFound
		}

		return fmt.Errorf("could not delete URL from database: %w", err)
	}

	return nil
}

func (s Service) GetRedirectionCount(ctx context.Context, short string) (int, error) {
	count, err := s.store.GetRedirectionCount(ctx, short)
	if err != nil {
		if err == ErrNotFound {
			return 0, ErrNotFound
		}

		return 0, fmt.Errorf("could not get URL redirection count from database: %w", err)
	}

	return count, nil
}
