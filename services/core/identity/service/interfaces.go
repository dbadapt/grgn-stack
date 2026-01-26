// Package service provides business logic for the identity domain.
package service

import (
	"context"

	"github.com/yourusername/grgn-stack/services/core/shared/generated/graphql/model"
)

// IUserService defines the contract for user business operations.
type IUserService interface {
	// GetCurrentUser retrieves the currently authenticated user.
	// Returns ErrNotAuthenticated if no user is in context.
	GetCurrentUser(ctx context.Context) (*model.User, error)

	// GetUserByID retrieves a user by their ID.
	// Returns ErrUserNotFound if the user doesn't exist.
	GetUserByID(ctx context.Context, id string) (*model.User, error)

	// UpdateProfile updates the current user's profile.
	// Returns ErrNotAuthenticated if no user is in context.
	UpdateProfile(ctx context.Context, input model.UpdateProfileInput) (*model.User, error)

	// DeleteAccount soft-deletes the current user's account.
	// Returns ErrNotAuthenticated if no user is in context.
	DeleteAccount(ctx context.Context) error

	// CreateUser creates a new user (internal use, e.g., seed command).
	// Returns ErrEmailTaken if the email already exists.
	CreateUser(ctx context.Context, email string, name *string) (*model.User, error)

	// GetUserByEmail retrieves a user by email (internal use).
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
}
