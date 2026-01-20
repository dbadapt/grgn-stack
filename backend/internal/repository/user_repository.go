package repository

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/yourusername/grgn-stack/internal/database"
)

// UserRepository handles User data operations
type UserRepository struct {
	db *database.Neo4jDB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *database.Neo4jDB) *UserRepository {
	return &UserRepository{db: db}
}

// User represents a user in the system
type User struct {
	ID        string
	Email     string
	Name      *string
	CreatedAt string
	UpdatedAt string
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, email string, name *string) (*User, error) {
	result, err := r.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			CREATE (u:User {
				id: randomUUID(),
				email: $email,
				name: $name,
				createdAt: datetime(),
				updatedAt: datetime()
			})
			RETURN u.id as id, u.email as email, u.name as name,
			       toString(u.createdAt) as createdAt, toString(u.updatedAt) as updatedAt
		`

		params := map[string]any{
			"email": email,
			"name":  name,
		}

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		record, err := result.Single(ctx)
		if err != nil {
			return nil, err
		}

		user := &User{
			ID:        record.Values[0].(string),
			Email:     record.Values[1].(string),
			CreatedAt: record.Values[3].(string),
			UpdatedAt: record.Values[4].(string),
		}

		if record.Values[2] != nil {
			nameStr := record.Values[2].(string)
			user.Name = &nameStr
		}

		return user, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return result.(*User), nil
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(ctx context.Context, id string) (*User, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (u:User {id: $id})
			RETURN u.id as id, u.email as email, u.name as name,
			       toString(u.createdAt) as createdAt, toString(u.updatedAt) as updatedAt
		`

		result, err := tx.Run(ctx, query, map[string]any{"id": id})
		if err != nil {
			return nil, err
		}

		if !result.Next(ctx) {
			return nil, fmt.Errorf("user not found")
		}

		record := result.Record()
		user := &User{
			ID:        record.Values[0].(string),
			Email:     record.Values[1].(string),
			CreatedAt: record.Values[3].(string),
			UpdatedAt: record.Values[4].(string),
		}

		if record.Values[2] != nil {
			nameStr := record.Values[2].(string)
			user.Name = &nameStr
		}

		return user, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return result.(*User), nil
}

// FindByEmail finds a user by email
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (u:User {email: $email})
			RETURN u.id as id, u.email as email, u.name as name,
			       toString(u.createdAt) as createdAt, toString(u.updatedAt) as updatedAt
		`

		result, err := tx.Run(ctx, query, map[string]any{"email": email})
		if err != nil {
			return nil, err
		}

		if !result.Next(ctx) {
			return nil, fmt.Errorf("user not found")
		}

		record := result.Record()
		user := &User{
			ID:        record.Values[0].(string),
			Email:     record.Values[1].(string),
			CreatedAt: record.Values[3].(string),
			UpdatedAt: record.Values[4].(string),
		}

		if record.Values[2] != nil {
			nameStr := record.Values[2].(string)
			user.Name = &nameStr
		}

		return user, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return result.(*User), nil
}

// Update updates a user's information
func (r *UserRepository) Update(ctx context.Context, id string, name *string) (*User, error) {
	result, err := r.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (u:User {id: $id})
			SET u.name = $name, u.updatedAt = datetime()
			RETURN u.id as id, u.email as email, u.name as name,
			       toString(u.createdAt) as createdAt, toString(u.updatedAt) as updatedAt
		`

		params := map[string]any{
			"id":   id,
			"name": name,
		}

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		if !result.Next(ctx) {
			return nil, fmt.Errorf("user not found")
		}

		record := result.Record()
		user := &User{
			ID:        record.Values[0].(string),
			Email:     record.Values[1].(string),
			CreatedAt: record.Values[3].(string),
			UpdatedAt: record.Values[4].(string),
		}

		if record.Values[2] != nil {
			nameStr := record.Values[2].(string)
			user.Name = &nameStr
		}

		return user, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return result.(*User), nil
}

// Delete deletes a user
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `MATCH (u:User {id: $id}) DETACH DELETE u`
		result, err := tx.Run(ctx, query, map[string]any{"id": id})
		if err != nil {
			return nil, err
		}

		summary, err := result.Consume(ctx)
		if err != nil {
			return nil, err
		}

		if summary.Counters().NodesDeleted() == 0 {
			return nil, fmt.Errorf("user not found")
		}

		return nil, nil
	})

	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
