package migrations

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// Migration001InitialSchema creates the initial database schema with constraints and indexes
var Migration001InitialSchema = Migration{
	Version:     1,
	Description: "Initial schema with User nodes, constraints, and indexes",

	Up: func(ctx context.Context, tx neo4j.ManagedTransaction) error {
		// Create unique constraint on User.email
		if _, err := tx.Run(ctx, `
			CREATE CONSTRAINT user_email_unique IF NOT EXISTS
			FOR (u:User) REQUIRE u.email IS UNIQUE
		`, nil); err != nil {
			return err
		}

		// Create unique constraint on User.id
		if _, err := tx.Run(ctx, `
			CREATE CONSTRAINT user_id_unique IF NOT EXISTS
			FOR (u:User) REQUIRE u.id IS UNIQUE
		`, nil); err != nil {
			return err
		}

		// Create index on User.createdAt for time-based queries
		if _, err := tx.Run(ctx, `
			CREATE INDEX user_created_at IF NOT EXISTS
			FOR (u:User) ON (u.createdAt)
		`, nil); err != nil {
			return err
		}

		return nil
	},

	Down: func(ctx context.Context, tx neo4j.ManagedTransaction) error {
		// Drop indexes
		indexes := []string{
			"DROP INDEX user_created_at IF EXISTS",
		}

		for _, query := range indexes {
			if _, err := tx.Run(ctx, query, nil); err != nil {
				return err
			}
		}

		// Drop constraints
		constraints := []string{
			"DROP CONSTRAINT user_email_unique IF EXISTS",
			"DROP CONSTRAINT user_id_unique IF EXISTS",
		}

		for _, query := range constraints {
			if _, err := tx.Run(ctx, query, nil); err != nil {
				return err
			}
		}

		return nil
	},
}
