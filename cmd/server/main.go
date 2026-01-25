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
	log.Println("Connecting to Neo4j database...")
	db, err := shared.NewNeo4jDB(cfg)
	if err != nil {
		log.Fatalf("Failed to create database connection: %v", err)
	}

	// Verify database connectivity
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.VerifyConnectivity(ctx); err != nil {
		log.Fatalf("Failed to connect to Neo4j: %v", err)
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
	log.Printf("Starting %s server on %s (environment: %s)", cfg.App.Name, addr, cfg.Server.Environment)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
