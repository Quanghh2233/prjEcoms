package main

import (
	"fmt"
	"log"

	"github.com/go-pg/pg/v10"
	"github.com/qhh/prjEcom/pkg/db/schema"
)

const (
	dbHost     = "localhost"
	dbPort     = "5432"
	dbUser     = "qhh"
	dbPassword = "2203"
	dbName     = "ecommerce"
)

func main() {
	// Connect to database
	db := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%s", dbHost, dbPort),
		User:     dbUser,
		Password: dbPassword,
		Database: dbName,
	})
	defer db.Close()

	// Initialize database schema
	err := schema.InitDatabase(db)
	if err != nil {
		log.Fatalf("Failed to initialize database schema: %v", err)
	}

	log.Println("Database schema initialized successfully")
}
