// Package repository provides data access for the tenant domain.
package repository

import (
	"context"

	"github.com/yourusername/grgn-stack/services/core/shared/generated/graphql/model"
)

// ITenantRepository defines the contract for tenant data access.
type ITenantRepository interface {
	// FindByID retrieves a tenant by their unique ID.
	// Returns ErrTenantNotFound if the tenant doesn't exist or is deleted.
	FindByID(ctx context.Context, id string) (*model.Tenant, error)

	// FindBySlug retrieves a tenant by their unique slug.
	// Returns ErrTenantNotFound if the tenant doesn't exist or is deleted.
	FindBySlug(ctx context.Context, slug string) (*model.Tenant, error)

	// FindByUserID retrieves all tenants a user is a member of.
	FindByUserID(ctx context.Context, userID string) ([]*model.Tenant, error)

	// Create creates a new tenant in the database.
	// Returns ErrSlugTaken if the slug already exists.
	Create(ctx context.Context, tenant *model.Tenant) (*model.Tenant, error)

	// Update updates an existing tenant.
	// Returns ErrTenantNotFound if the tenant doesn't exist.
	Update(ctx context.Context, id string, input model.UpdateTenantInput) (*model.Tenant, error)

	// Delete soft-deletes a tenant by setting their status to DELETED.
	// Returns ErrTenantNotFound if the tenant doesn't exist.
	Delete(ctx context.Context, id string) error

	// ExistsBySlug checks if a tenant with the given slug exists.
	ExistsBySlug(ctx context.Context, slug string) (bool, error)

	// GetMemberCount returns the number of members in a tenant.
	GetMemberCount(ctx context.Context, tenantID string) (int, error)
}

// IMembershipRepository defines the contract for membership data access.
type IMembershipRepository interface {
	// FindByID retrieves a membership by its unique ID.
	// Returns ErrMembershipNotFound if the membership doesn't exist.
	FindByID(ctx context.Context, id string) (*model.Membership, error)

	// FindByTenantID retrieves all memberships for a tenant.
	FindByTenantID(ctx context.Context, tenantID string) ([]*model.Membership, error)

	// FindByUserID retrieves all memberships for a user.
	FindByUserID(ctx context.Context, userID string) ([]*model.Membership, error)

	// FindByUserAndTenant retrieves a membership by user and tenant.
	// Returns ErrMembershipNotFound if the membership doesn't exist.
	FindByUserAndTenant(ctx context.Context, userID, tenantID string) (*model.Membership, error)

	// Create creates a new membership.
	// Returns ErrAlreadyMember if the user is already a member.
	Create(ctx context.Context, userID, tenantID string, role model.MembershipRole, invitedByID *string) (*model.Membership, error)

	// UpdateRole updates a membership's role.
	// Returns ErrMembershipNotFound if the membership doesn't exist.
	UpdateRole(ctx context.Context, id string, role model.MembershipRole) (*model.Membership, error)

	// Delete removes a membership.
	// Returns ErrMembershipNotFound if the membership doesn't exist.
	Delete(ctx context.Context, id string) error

	// CountOwners returns the number of owners in a tenant.
	CountOwners(ctx context.Context, tenantID string) (int, error)

	// GetTenantIDByMembershipID returns the tenant ID for a membership.
	GetTenantIDByMembershipID(ctx context.Context, membershipID string) (string, error)

	// GetUserIDByMembershipID returns the user ID for a membership.
	GetUserIDByMembershipID(ctx context.Context, membershipID string) (string, error)
}
