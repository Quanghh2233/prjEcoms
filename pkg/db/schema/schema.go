package schema

import (
	"fmt"
	"log"

	"github.com/go-pg/pg/v10"
	"github.com/qhh/prjEcom/pkg/models"
)

// InitDatabase initializes the database schema
func InitDatabase(db *pg.DB) error {
	// Create enum types
	err := createEnumTypes(db)
	if err != nil {
		return fmt.Errorf("failed to create enum types: %w", err)
	}

	// Create tables
	err = models.CreateSchema(db)
	if err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	log.Println("Database schema initialized successfully")
	return nil
}

// createEnumTypes creates the necessary enum types in PostgreSQL
func createEnumTypes(db *pg.DB) error {
	// Create user_role enum if it doesn't exist
	_, err := db.Exec(`DO $$ 
	BEGIN 
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role') THEN
			CREATE TYPE user_role AS ENUM ('admin', 'seller', 'buyer');
		END IF;
	END $$;`)
	if err != nil {
		return err
	}

	// Create order_status enum if it doesn't exist
	_, err = db.Exec(`DO $$ 
	BEGIN 
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'order_status') THEN
			CREATE TYPE order_status AS ENUM ('pending', 'paid', 'shipped', 'delivered', 'canceled');
		END IF;
	END $$;`)

	return err
}
