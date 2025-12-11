package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Master wraps a pgx connection pool for the master database.
type Master struct {
	Pool *pgxpool.Pool
}

// NewMaster connects to the master database.
func NewMaster(dsn string) (*Master, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse master dsn: %w", err)
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, fmt.Errorf("connect master: %w", err)
	}
	return &Master{Pool: pool}, nil
}

// Close closes the underlying pool.
func (m *Master) Close() {
	m.Pool.Close()
}

// EnsureSchema creates the articles table if it does not exist.
// This is a safeguard in addition to the SQL init script.
func EnsureSchema(m *Master) error {
	_, err := m.Pool.Exec(context.Background(), `
CREATE TABLE IF NOT EXISTS articles (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    author TEXT NOT NULL,
    region TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
`)
	return err
}



