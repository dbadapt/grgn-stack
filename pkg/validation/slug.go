// Package validation provides input validation utilities.
package validation

import (
	"regexp"

	"github.com/yourusername/grgn-stack/pkg/errors"
)

// slugRegex matches valid slug format: letters, numbers, hyphens, underscores
// Length: 3-50 characters
var slugRegex = regexp.MustCompile(`^[A-Za-z0-9_-]{3,50}$`)

// ValidateSlug checks if slug matches allowed format.
// Allowed: A-Z, a-z, 0-9, hyphen (-), underscore (_)
// Length: 3-50 characters
func ValidateSlug(slug string) error {
	if !slugRegex.MatchString(slug) {
		return errors.NewValidationError("slug",
			"must be 3-50 characters, containing only letters, numbers, hyphens, and underscores")
	}
	return nil
}

// IsValidSlug returns true if the slug is valid
func IsValidSlug(slug string) bool {
	return slugRegex.MatchString(slug)
}
