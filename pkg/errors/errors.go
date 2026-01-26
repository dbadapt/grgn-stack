// Package errors provides custom error types for the GRGN stack.
package errors

import "errors"

// Sentinel errors for common cases
var (
	// Not found errors
	ErrNotFound           = errors.New("resource not found")
	ErrUserNotFound       = errors.New("user not found")
	ErrTenantNotFound     = errors.New("tenant not found")
	ErrMembershipNotFound = errors.New("membership not found")

	// Auth errors
	ErrNotAuthenticated = errors.New("user not authenticated")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrForbidden        = errors.New("forbidden: insufficient permissions")

	// Validation errors
	ErrInvalidInput = errors.New("invalid input")
	ErrInvalidSlug  = errors.New("invalid slug format")
	ErrSlugTaken    = errors.New("slug already taken")
	ErrEmailTaken   = errors.New("email already taken")

	// Business rule errors
	ErrLastOwner     = errors.New("cannot remove or demote the last owner")
	ErrAlreadyMember = errors.New("user is already a member")
	ErrNotMember     = errors.New("user is not a member of this tenant")
	ErrCannotLeave   = errors.New("cannot leave: you are the last owner")
)

// ValidationError wraps validation errors with field info
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{Field: field, Message: message}
}

// Is checks if target error matches
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As finds the first error in err's chain that matches target
func As(err error, target any) bool {
	return errors.As(err, target)
}

// Wrap wraps an error with additional context
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return errors.New(message + ": " + err.Error())
}
