package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// Open opens and pings a *sql.DB for the given Postgres DSN.
func Open(dbpath string) (*sql.DB, error) {
	conn, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}
	return conn, nil
}
