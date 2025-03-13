package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/qhh/prjEcom/pkg/api/routes"
	"github.com/qhh/prjEcom/pkg/config"
	"github.com/qhh/prjEcom/pkg/db"
	dbinit "github.com/qhh/prjEcom/pkg/db/dbinit"
	"github.com/qhh/prjEcom/pkg/db/store"
	"github.com/qhh/prjEcom/pkg/utils"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	pgDB, err := db.Connect(&cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pgDB.Close()

	// Initialize database if needed
	if os.Getenv("INIT_DB") == "true" {
		if err := dbinit.InitializeDatabase(&cfg); err != nil {
			log.Printf("Warning: failed to initialize database: %v", err)
		} else {
			log.Println("Database initialized successfully")
		}
	}

	// Create store
	store := store.NewStore(pgDB)

	// Create JWT maker
	jwtMaker, err := utils.NewJWTMaker(cfg.JWTSecret)
	if err != nil {
		log.Fatalf("Failed to create JWT maker: %v", err)
	}

	// Setup router
	router := routes.SetupRouter(store, jwtMaker)

	// Start server
	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server starting on %s", serverAddr)

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := router.Run(serverAddr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")
}
