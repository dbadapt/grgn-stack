package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/yourusername/grgn-stack/pkg/config"
	shared "github.com/yourusername/grgn-stack/services/core/shared/controller"
	"github.com/yourusername/grgn-stack/services/core/shared/generated/graphql"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize Neo4j database connection
	log.Printf("Connecting to Neo4j database at %s...", cfg.Database.Neo4jURI)
	db, err := shared.NewNeo4jDB(cfg)
	if err != nil {
		log.Fatalf("Failed to create database connection: %v", err)
	}

	// Verify database connectivity with retry
	var dbConnected bool
	for i := 0; i < 10; i++ {
		log.Printf("Verifying database connectivity (attempt %d/10)...", i+1)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err = db.VerifyConnectivity(ctx)
		cancel()

		if err == nil {
			dbConnected = true
			break
		}

		log.Printf("Database not ready: %v. Retrying in 2 seconds...", err)
		time.Sleep(2 * time.Second)
	}

	if !dbConnected {
		log.Fatalf("Failed to connect to Neo4j after 10 attempts: %v", err)
	}
	log.Println("Successfully connected to Neo4j")

	// Set up graceful shutdown
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-shutdownChan
		log.Println("Shutting down gracefully...")

		// Close database connection
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := db.Close(ctx); err != nil {
			log.Printf("Error closing database: %v", err)
		} else {
			log.Println("Database connection closed")
		}

		os.Exit(0)
	}()

	// Set Gin mode based on environment
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else if cfg.IsStaging() {
		gin.SetMode(gin.TestMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.Default()

	// Create ping handler and register route
	pingHandler := shared.NewPingHandler(db, cfg)
	r.GET("/ping", pingHandler.HandlePing)

	// GraphQL setup
	gqlResolver := &graphql.Resolver{}
	gqlServer := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: gqlResolver}))

	// GraphQL endpoints
	r.POST("/graphql", func(c *gin.Context) {
		gqlServer.ServeHTTP(c.Writer, c.Request)
	})

	// GraphQL Playground (only in development)
	if !cfg.IsProduction() {
		r.GET("/graphql", func(c *gin.Context) {
			playground.Handler("GRGN Stack GraphQL Playground", "/graphql").ServeHTTP(c.Writer, c.Request)
		})
		log.Printf("GraphQL Playground available at http://%s:%s/graphql", cfg.Server.Host, cfg.Server.Port)
	}

	// Start server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting %s server...", cfg.App.Name)
	log.Printf("Environment: %s", cfg.Server.Environment)
	log.Printf("Version: %s", cfg.App.Version)
	log.Printf("Listening on: http://%s", addr)
	log.Printf("GraphQL endpoint: http://%s/graphql", addr)
	if !cfg.IsProduction() {
		log.Printf("GraphQL Playground: http://%s/graphql", addr)
	}

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
