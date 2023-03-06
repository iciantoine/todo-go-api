package database

import (
	"database/sql"
	"fmt"
	"time"
)

// Connect connects to a database and sets default options.
func Connect(driver, dsn string) (*sql.DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("could not open DB connection pool: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("could not ping DB: %w", err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	return db, nil
}
