package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// Connect establishes a connection to the PostgreSQL database using the provided DSN.
// It returns a *sql.DB instance or an error if the connection fails.
func Connect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Connection Pooling
	if err := configurePooling(db); err != nil {
		log.Printf("Warning: Failed to configure connection pooling: %v", err)
	}

	log.Println("Connected to database successfully")
	return db, nil
}

func configurePooling(db *sql.DB) error {
	defaultMaxOpen := 25
	defaultMaxIdle := 25
	defaultMaxLifetime := 5 * time.Minute

	if maxOpen := os.Getenv("DB_MAX_OPEN_CONNS"); maxOpen != "" {
		if val, err := strconv.Atoi(maxOpen); err == nil {
			db.SetMaxOpenConns(val)
		}
	} else {
		db.SetMaxOpenConns(defaultMaxOpen)
	}

	if maxIdle := os.Getenv("DB_MAX_IDLE_CONNS"); maxIdle != "" {
		if val, err := strconv.Atoi(maxIdle); err == nil {
			db.SetMaxIdleConns(val)
		}
	} else {
		db.SetMaxIdleConns(defaultMaxIdle)
	}

	if lifetime := os.Getenv("DB_CONN_MAX_LIFETIME"); lifetime != "" {
		if val, err := time.ParseDuration(lifetime); err == nil {
			db.SetConnMaxLifetime(val)
		}
	} else {
		db.SetConnMaxLifetime(defaultMaxLifetime)
	}

	return nil
}

// Migrate runs database migrations using golang-migrate.
func Migrate(db *sql.DB, databaseName, migrationsPath string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		databaseName,
		driver,
	)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}

	log.Println("Database migration completed successfully")
	return nil
}
