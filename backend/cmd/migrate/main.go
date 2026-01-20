package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/yourusername/grgn-stack/internal/database"
	"github.com/yourusername/grgn-stack/internal/database/migrations"
	"github.com/yourusername/grgn-stack/pkg/config"
)

func main() {
	// Define command line flags
	command := flag.String("command", "status", "Migration command: up, down, status")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Connecting to Neo4j at %s (environment: %s)", cfg.Database.Neo4jURI, cfg.Server.Environment)

	// Initialize database connection
	db, err := database.NewNeo4jDB(cfg)
	if err != nil {
		log.Fatalf("Failed to create database connection: %v", err)
	}

	// Verify connectivity
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.VerifyConnectivity(ctx); err != nil {
		log.Fatalf("Failed to connect to Neo4j: %v", err)
	}

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := db.Close(ctx); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	// Create migrator with all migrations
	migrator := migrations.NewMigratorWithAll(db)

	// Execute command
	ctx = context.Background()

	switch *command {
	case "up":
		log.Println("Running migrations...")
		if err := migrator.Up(ctx); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		log.Println("✓ Migrations completed successfully")

	case "down":
		log.Println("Rolling back last migration...")
		if err := migrator.Down(ctx); err != nil {
			log.Fatalf("Rollback failed: %v", err)
		}
		log.Println("✓ Rollback completed successfully")

	case "status":
		if err := migrator.Status(ctx); err != nil {
			log.Fatalf("Failed to get status: %v", err)
		}

	default:
		fmt.Printf("Unknown command: %s\n", *command)
		fmt.Println("Available commands: up, down, status")
		os.Exit(1)
	}
}
