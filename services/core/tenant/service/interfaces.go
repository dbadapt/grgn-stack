// Package service provides business logic for the tenant domain.
package service

import (
	"context"

	"github.com/yourusername/grgn-stack/services/core/shared/generated/graphql/model"
)

// ITenantService defines the contract for tenant business operations.
type ITenantService interface {
	// Tenant operations

	// GetTenant retrieves a tenant by ID.
	GetTenant(ctx context.Context, id string) (*model.Tenant, error)

	// GetTenantBySlug retrieves a tenant by slug.
	GetTenantBySlug(ctx context.Context, slug string) (*model.Tenant, error)

	// GetMyTenants retrieves all tenants the current user is a member of.
	GetMyTenants(ctx context.Context) ([]*model.Tenant, error)

	// CreateTenant creates a new tenant with the current user as owner.
	CreateTenant(ctx context.Context, input model.CreateTenantInput) (*model.Tenant, error)

	// UpdateTenant updates a tenant. Requires ADMIN+ role.
	UpdateTenant(ctx context.Context, id string, input model.UpdateTenantInput) (*model.Tenant, error)

	// DeleteTenant soft-deletes a tenant. Requires OWNER role.
	DeleteTenant(ctx context.Context, id string) (bool, error)

	// Membership operations

	// GetTenantMembers retrieves all members of a tenant.
	GetTenantMembers(ctx context.Context, tenantID string) ([]*model.Membership, error)

	// InviteMember invites a user to a tenant. Requires ADMIN+ role.
	InviteMember(ctx context.Context, tenantID string, input model.InviteMemberInput) (*model.Membership, error)

	// UpdateMemberRole updates a member's role. Requires OWNER role.
	UpdateMemberRole(ctx context.Context, membershipID string, role model.MembershipRole) (*model.Membership, error)

	// RemoveMember removes a member from a tenant. Requires ADMIN+ role.
	RemoveMember(ctx context.Context, membershipID string) (bool, error)

	// LeaveTenant removes the current user from a tenant.
	LeaveTenant(ctx context.Context, tenantID string) (bool, error)
}
