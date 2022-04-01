package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/nerock/urlshort/url"
)

type URLStore struct {
	db *sql.DB
}

const (
	createURLTable = `CREATE TABLE IF NOT EXISTS url (short TEXT NOT NULL PRIMARY KEY, long TEXT NOT NULL, count INTEGER DEFAULT 0)`

	createURL = `INSERT INTO url (short, long) VALUES (?, ?)`
	getURL    = `SELECT long FROM url WHERE short = ?`
	deleteURL = `DELETE FROM url WHERE short = ?`

	incrementRedirectionCount = `UPDATE url SET count = count + 1 WHERE short = ?`
	getRedirectiontCount      = `SELECT count FROM url WHERE short = ?`
)

// NewURLStore instantiates a new url store with sqlite
func NewURLStore(db *sql.DB) (URLStore, error) {
	if _, err := db.Exec(createURLTable); err != nil {
		return URLStore{}, fmt.Errorf("could not create url table: %w", err)
	}

	return URLStore{db}, nil
}

// AddURL saves a new url
func (u URLStore) AddURL(ctx context.Context, short, long string) error {
	if _, err := u.db.ExecContext(ctx, createURL, short, long); err != nil {
		return fmt.Errorf("save url in database: %w", err)
	}

	return nil
}

// GetURL gets a long url from the id
func (u URLStore) GetURL(ctx context.Context, short string) (string, error) {
	row := u.db.QueryRowContext(ctx, getURL, short)
	if row.Err() != nil {
		return "", fmt.Errorf("get url from database: %w", row.Err())
	}

	var long string
	if err := row.Scan(&long); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", url.ErrNotFound
		}

		return "", fmt.Errorf("parse url from database: %w", err)
	}

	return long, nil
}

// DeleteURL deletes an url from the id
func (u URLStore) DeleteURL(ctx context.Context, short string) error {
	if _, err := u.db.ExecContext(ctx, deleteURL, short); err != nil {
		return fmt.Errorf("delete url from database: %w", err)
	}

	return nil
}

// IncrementRedirectionCount increments the count by one
func (u URLStore) IncrementRedirectionCount(ctx context.Context, short string) error {
	if _, err := u.db.ExecContext(ctx, incrementRedirectionCount, short); err != nil {
		return fmt.Errorf("save url in database: %w", err)
	}

	return nil
}

// GetRedirectionCount gets the count of a url
func (u URLStore) GetRedirectionCount(ctx context.Context, short string) (int, error) {
	row := u.db.QueryRowContext(ctx, getRedirectiontCount, short)
	if row.Err() != nil {
		return 0, fmt.Errorf("get redirection count from database: %w", row.Err())
	}

	var count int
	if err := row.Scan(&count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, url.ErrNotFound
		}

		return 0, fmt.Errorf("parse count from database: %w", err)
	}

	return count, nil
}
