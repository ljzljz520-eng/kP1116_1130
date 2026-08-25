package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

func Open(path string) (*Store, error) {
	if strings.TrimSpace(path) == "" {
		return nil, fmt.Errorf("database path is empty")
	}
	db, err := sql.Open("sqlite", filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	store := &Store{db: db}
	if err := store.migrate(context.Background()); err != nil {
		db.Close()
		return nil, err
	}
	return store, nil
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *Store) Ping(ctx context.Context) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("store is closed")
	}
	return s.db.PingContext(ctx)
}

// SetMaxOpenConns forwards to the underlying connection pool. It lets callers
// (notably tests backed by a shared in-memory database) pin the pool to a
// single connection so a schema created on one connection stays visible to
// every subsequent query.
func (s *Store) SetMaxOpenConns(n int) {
	if s == nil || s.db == nil {
		return
	}
	s.db.SetMaxOpenConns(n)
}

func (s *Store) migrate(ctx context.Context) error {
	for _, statement := range schemaStatements() {
		if _, err := s.db.ExecContext(ctx, statement); err != nil {
			return fmt.Errorf("migrate schema: %w", err)
		}
	}
	return nil
}
