package database

import (
	"database/sql"
	_ "embed"

	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schema string

type DB struct {
	*sql.DB
}

func New(dsn string) (*DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)

	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}

	// Run migrations
	if err := runMigrations(db); err != nil {
		return nil, err
	}

	return &DB{DB: db}, nil
}

func runMigrations(db *sql.DB) error {
	// Check if default_branch column exists in repositories table
	var columnExists bool
	row := db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('repositories') WHERE name='default_branch'")
	if err := row.Scan(&columnExists); err != nil {
		return err
	}

	// Add default_branch column if it doesn't exist
	if !columnExists {
		_, err := db.Exec("ALTER TABLE repositories ADD COLUMN default_branch TEXT NOT NULL DEFAULT 'main'")
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}
