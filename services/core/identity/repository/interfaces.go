// Package repository provides data access for the identity domain.
package repository

import (
	"context"

	"github.com/yourusername/grgn-stack/services/core/shared/generated/graphql/model"
)

// IUserRepository defines the contract for user data access.
type IUserRepository interface {
	// FindByID retrieves a user by their unique ID.
	// Returns ErrUserNotFound if the user doesn't exist or is deleted.
	FindByID(ctx context.Context, id string) (*model.User, error)

	// FindByEmail retrieves a user by their email address.
	// Returns ErrUserNotFound if the user doesn't exist or is deleted.
	FindByEmail(ctx context.Context, email string) (*model.User, error)

	// Create creates a new user in the database.
	// Returns ErrEmailTaken if the email already exists.
	Create(ctx context.Context, user *model.User) (*model.User, error)

	// Update updates an existing user's profile.
	// Returns ErrUserNotFound if the user doesn't exist.
	Update(ctx context.Context, id string, input model.UpdateProfileInput) (*model.User, error)

	// Delete soft-deletes a user by setting their status to DELETED.
	// Returns ErrUserNotFound if the user doesn't exist.
	Delete(ctx context.Context, id string) error

	// List retrieves users with pagination.
	List(ctx context.Context, limit, offset int) ([]*model.User, error)

	// ExistsByEmail checks if a user with the given email exists.
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
