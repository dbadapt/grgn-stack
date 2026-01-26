// Package auth provides authentication context helpers.
// This is a stub that will be replaced when core/auth is implemented.
package auth

import (
	"context"

	"github.com/yourusername/grgn-stack/pkg/errors"
)

type contextKey string

// UserIDKey is the context key for storing user ID
const UserIDKey contextKey = "userID"

// GetUserID extracts the user ID from context.
// Returns ErrNotAuthenticated if no user ID is present.
func GetUserID(ctx context.Context) (string, error) {
	id, ok := ctx.Value(UserIDKey).(string)
	if !ok || id == "" {
		return "", errors.ErrNotAuthenticated
	}
	return id, nil
}

// WithUserID adds user ID to context
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// MustGetUserID extracts the user ID from context or panics.
// Use only when you're certain the user is authenticated.
func MustGetUserID(ctx context.Context) string {
	id, err := GetUserID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}
