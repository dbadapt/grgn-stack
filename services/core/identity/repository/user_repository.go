package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/yourusername/grgn-stack/pkg/errors"
	shared "github.com/yourusername/grgn-stack/services/core/shared/controller"
	"github.com/yourusername/grgn-stack/services/core/shared/generated/graphql/model"
)

// UserRepository implements IUserRepository using Neo4j.
type UserRepository struct {
	db shared.IDatabase
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(db shared.IDatabase) *UserRepository {
	return &UserRepository{db: db}
}

// FindByID retrieves a user by their unique ID.
func (r *UserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (u:User {id: $id})
			WHERE u.status <> 'DELETED'
			RETURN u
		`, map[string]any{"id": id})
		if err != nil {
			return nil, err
		}

		record, err := result.Single(ctx)
		if err != nil {
			return nil, errors.ErrUserNotFound
		}

		return r.mapRecordToUser(record, "u")
	})
	if err != nil {
		return nil, err
	}
	return result.(*model.User), nil
}

// FindByEmail retrieves a user by their email address.
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (u:User {email: $email})
			WHERE u.status <> 'DELETED'
			RETURN u
		`, map[string]any{"email": email})
		if err != nil {
			return nil, err
		}

		record, err := result.Single(ctx)
		if err != nil {
			return nil, errors.ErrUserNotFound
		}

		return r.mapRecordToUser(record, "u")
	})
	if err != nil {
		return nil, err
	}
	return result.(*model.User), nil
}

// Create creates a new user in the database.
func (r *UserRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	// Generate ID if not provided
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	if user.Status == "" {
		user.Status = model.UserStatusActive
	}

	result, err := r.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		// Check if email already exists
		checkResult, err := tx.Run(ctx, `
			MATCH (u:User {email: $email})
			RETURN count(u) > 0 as exists
		`, map[string]any{"email": user.Email})
		if err != nil {
			return nil, err
		}

		checkRecord, err := checkResult.Single(ctx)
		if err != nil {
			return nil, err
		}

		if exists, _ := checkRecord.Get("exists"); exists.(bool) {
			return nil, errors.ErrEmailTaken
		}

		// Create the user
		params := map[string]any{
			"id":        user.ID,
			"email":     user.Email,
			"name":      user.Name,
			"avatarUrl": user.AvatarURL,
			"status":    string(user.Status),
		}

		result, err := tx.Run(ctx, `
			CREATE (u:User {
				id: $id,
				email: $email,
				name: $name,
				avatarUrl: $avatarUrl,
				status: $status,
				createdAt: datetime(),
				updatedAt: datetime()
			})
			RETURN u
		`, params)
		if err != nil {
			return nil, err
		}

		record, err := result.Single(ctx)
		if err != nil {
			return nil, err
		}

		return r.mapRecordToUser(record, "u")
	})
	if err != nil {
		return nil, err
	}
	return result.(*model.User), nil
}

// Update updates an existing user's profile.
func (r *UserRepository) Update(ctx context.Context, id string, input model.UpdateProfileInput) (*model.User, error) {
	result, err := r.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		params := map[string]any{
			"id":        id,
			"updatedAt": time.Now(),
		}

		// Build SET clause dynamically
		setClause := "u.updatedAt = datetime($updatedAt)"
		if input.Name != nil {
			params["name"] = *input.Name
			setClause += ", u.name = $name"
		}
		if input.AvatarURL != nil {
			params["avatarUrl"] = *input.AvatarURL
			setClause += ", u.avatarUrl = $avatarUrl"
		}

		query := `
			MATCH (u:User {id: $id})
			WHERE u.status <> 'DELETED'
			SET ` + setClause + `
			RETURN u
		`

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		record, err := result.Single(ctx)
		if err != nil {
			return nil, errors.ErrUserNotFound
		}

		return r.mapRecordToUser(record, "u")
	})
	if err != nil {
		return nil, err
	}
	return result.(*model.User), nil
}

// Delete soft-deletes a user by setting their status to DELETED.
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (u:User {id: $id})
			WHERE u.status <> 'DELETED'
			SET u.status = 'DELETED', u.updatedAt = datetime()
			RETURN u
		`, map[string]any{"id": id})
		if err != nil {
			return nil, err
		}

		_, err = result.Single(ctx)
		if err != nil {
			return nil, errors.ErrUserNotFound
		}

		return nil, nil
	})
	return err
}

// List retrieves users with pagination.
func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*model.User, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (u:User)
			WHERE u.status <> 'DELETED'
			RETURN u
			ORDER BY u.createdAt DESC
			SKIP $offset
			LIMIT $limit
		`, map[string]any{"limit": limit, "offset": offset})
		if err != nil {
			return nil, err
		}

		var users []*model.User
		for result.Next(ctx) {
			user, err := r.mapRecordToUser(result.Record(), "u")
			if err != nil {
				return nil, err
			}
			users = append(users, user)
		}

		return users, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*model.User), nil
}

// ExistsByEmail checks if a user with the given email exists.
func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (u:User {email: $email})
			WHERE u.status <> 'DELETED'
			RETURN count(u) > 0 as exists
		`, map[string]any{"email": email})
		if err != nil {
			return nil, err
		}

		record, err := result.Single(ctx)
		if err != nil {
			return false, nil
		}

		exists, _ := record.Get("exists")
		return exists.(bool), nil
	})
	if err != nil {
		return false, err
	}
	return result.(bool), nil
}

// mapRecordToUser converts a Neo4j record to a User model.
func (r *UserRepository) mapRecordToUser(record *neo4j.Record, key string) (*model.User, error) {
	nodeVal, ok := record.Get(key)
	if !ok {
		return nil, errors.ErrUserNotFound
	}

	node := nodeVal.(neo4j.Node)
	props := node.Props

	user := &model.User{
		ID:     props["id"].(string),
		Email:  props["email"].(string),
		Status: model.UserStatus(props["status"].(string)),
	}

	if name, ok := props["name"]; ok && name != nil {
		nameStr := name.(string)
		user.Name = &nameStr
	}

	if avatarURL, ok := props["avatarUrl"]; ok && avatarURL != nil {
		avatarStr := avatarURL.(string)
		user.AvatarURL = &avatarStr
	}

	if createdAt, ok := props["createdAt"]; ok {
		user.CreatedAt = createdAt.(time.Time)
	}

	if updatedAt, ok := props["updatedAt"]; ok {
		user.UpdatedAt = updatedAt.(time.Time)
	}

	return user, nil
}

// Ensure UserRepository implements IUserRepository
var _ IUserRepository = (*UserRepository)(nil)
