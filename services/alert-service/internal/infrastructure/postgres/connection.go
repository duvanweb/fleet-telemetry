package postgres

import (
	"context"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"fleet/alert-service/internal/core/ports/repositories"
)

// Configuration holds PostgreSQL connection settings.
type Configuration struct {
	Host     string `env:"DB_HOST"     envDefault:"localhost"`
	Port     string `env:"DB_PORT"     envDefault:"5432"`
	User     string `env:"DB_USER"     envDefault:"postgres"`
	Password string `env:"DB_PASSWORD" envDefault:"postgres"`
	DBName   string `env:"DB_NAME"     envDefault:"alert_db"`
	SSLMode  string `env:"DB_SSLMODE"  envDefault:"disable"`
}

func (c *Configuration) dsn() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

// NewConnection opens a PostgreSQL connection pool and verifies connectivity.
func NewConnection(cfg *Configuration) (repositories.Databaser, error) {
	db, err := openDB("postgres", cfg.dsn())
	if err != nil {
		return nil, fmt.Errorf("opening postgres connection: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(2 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("pinging postgres: %w", err)
	}

	return db, nil
}
