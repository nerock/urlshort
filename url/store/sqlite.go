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
	createURLTable = `CREATE TABLE IF NOT EXISTS url (short TEXT NOT NULL PRIMARY KEY, long TEXT NOT NULL)`
	createURL      = `INSERT INTO url (short, long) VALUES (?, ?)`
	getURL         = `SELECT long FROM url WHERE short = ?`
	deleteURL      = `DELETE FROM url WHERE short = ?`
)

func NewURLStore(db *sql.DB) (URLStore, error) {
	if _, err := db.Exec(createURLTable); err != nil {
		return URLStore{}, fmt.Errorf("could not create url table: %w", err)
	}

	return URLStore{db}, nil
}

func (u URLStore) AddURL(ctx context.Context, short, long string) error {
	if _, err := u.db.ExecContext(ctx, createURL, short, long); err != nil {
		return fmt.Errorf("save url in database: %w", err)
	}

	return nil
}

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

func (u URLStore) DeleteURL(ctx context.Context, short string) error {
	if _, err := u.db.ExecContext(ctx, deleteURL, short); err != nil {
		return fmt.Errorf("delete url from database: %w", err)
	}

	return nil
}
