package sqlite

import (
	"database/sql"
	"fmt"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)

	if err != nil {
		return nil, fmt.Errorf("Can't open database file %s: %w", path, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Database ping failed: %w", err)
	}

	return &Storage{db}, nil
}
