package db

import (
	"context"
	"fmt"
	"net/url"

	"github.com/go-pg/pg/v10"
	"github.com/qhh/prjEcom/pkg/config"
)

// Connect establishes a connection to the PostgreSQL database
func Connect(cfg *config.Config) (*pg.DB, error) {
	// Parse the connection string to extract individual components
	parsedURL, err := url.Parse(cfg.DBSource)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DB URL: %w", err)
	}

	password, _ := parsedURL.User.Password()
	dbname := parsedURL.Path[1:] // Remove leading '/'

	opt := &pg.Options{
		Addr:     parsedURL.Host,
		User:     parsedURL.User.Username(),
		Password: password,
		Database: dbname,
	}

	db := pg.Connect(opt)

	// Check if the connection is working
	err = db.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
