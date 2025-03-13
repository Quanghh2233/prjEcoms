package dbinit

import (
	"fmt"
	"log"

	"github.com/go-pg/pg/v10"
	"github.com/qhh/prjEcom/pkg/config"
	"github.com/qhh/prjEcom/pkg/db/schema"
)

// InitializeDatabase sets up the database schema and initial data
func InitializeDatabase(cfg *config.Config) error {
	// Connect to DB
	opt, err := pg.ParseURL(cfg.DBSource)
	if err != nil {
		return fmt.Errorf("failed to parse DB URL: %w", err)
	}

	db := pg.Connect(opt)
	defer db.Close()

	// Create schema
	if err := schema.InitDatabase(db); err != nil {
		return fmt.Errorf("failed to initialize database schema: %w", err)
	}

	log.Println("Database initialized successfully")
	return nil
}
