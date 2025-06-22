package database

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func NewPostgresDB(databaseUrl string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", databaseUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetConnMaxIdleTime(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Could not ping database: %v", err)
	}

	log.Println("Successfully connected to PostgreSQL!")

	return db, nil
}
