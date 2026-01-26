package service

import (
	"context"

	"github.com/yourusername/grgn-stack/pkg/auth"
	"github.com/yourusername/grgn-stack/pkg/errors"
	"github.com/yourusername/grgn-stack/pkg/validation"
	identityRepo "github.com/yourusername/grgn-stack/services/core/identity/repository"
	"github.com/yourusername/grgn-stack/services/core/shared/generated/graphql/model"
	"github.com/yourusername/grgn-stack/services/core/tenant/repository"
)

// TenantService implements ITenantService with business logic.
type TenantService struct {
	tenantRepo     repository.ITenantRepository
	membershipRepo repository.IMembershipRepository
	userRepo       identityRepo.IUserRepository
}

// NewTenantService creates a new TenantService.
func NewTenantService(
	tenantRepo repository.ITenantRepository,
	membershipRepo repository.IMembershipRepository,
	userRepo identityRepo.IUserRepository,
) *TenantService {
	return &TenantService{
		tenantRepo:     tenantRepo,
		membershipRepo: membershipRepo,
		userRepo:       userRepo,
	}
}

// Role hierarchy: OWNER > ADMIN > MEMBER > VIEWER
var roleOrder = map[model.MembershipRole]int{
	model.MembershipRoleViewer: 1,
	model.MembershipRoleMember: 2,
	model.MembershipRoleAdmin:  3,
	model.MembershipRoleOwner:  4,
}

// hasMinRole checks if the actual role meets the minimum required role.
func hasMinRole(actual, required model.MembershipRole) bool {
	return roleOrder[actual] >= roleOrder[required]
}

// requireRole checks if the current user has at least the required role in a tenant.
func (s *TenantService) requireRole(ctx context.Context, tenantID string, minRole model.MembershipRole) (*model.Membership, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	membership, err := s.membershipRepo.FindByUserAndTenant(ctx, userID, tenantID)
	if err != nil {
		return nil, errors.ErrNotMember
	}

	if !hasMinRole(membership.Role, minRole) {
		return nil, errors.ErrForbidden
	}

	return membership, nil
}

// GetTenant retrieves a tenant by ID.
func (s *TenantService) GetTenant(ctx context.Context, id string) (*model.Tenant, error) {
	return s.tenantRepo.FindByID(ctx, id)
}

// GetTenantBySlug retrieves a tenant by slug.
func (s *TenantService) GetTenantBySlug(ctx context.Context, slug string) (*model.Tenant, error) {
	return s.tenantRepo.FindBySlug(ctx, slug)
}

// GetMyTenants retrieves all tenants the current user is a member of.
func (s *TenantService) GetMyTenants(ctx context.Context) ([]*model.Tenant, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	return s.tenantRepo.FindByUserID(ctx, userID)
}

// CreateTenant creates a new tenant with the current user as owner.
func (s *TenantService) CreateTenant(ctx context.Context, input model.CreateTenantInput) (*model.Tenant, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	// Validate slug
	if err := validation.ValidateSlug(input.Slug); err != nil {
		return nil, errors.ErrInvalidSlug
	}

	// Set default plan if not provided
	plan := model.TenantPlanFree
	if input.Plan != nil {
		plan = *input.Plan
	}

	// Create tenant
	tenant := &model.Tenant{
		Name:          input.Name,
		Slug:          input.Slug,
		Plan:          plan,
		Status:        model.TenantStatusActive,
		IsolationMode: model.TenantIsolationModeShared,
	}

	createdTenant, err := s.tenantRepo.Create(ctx, tenant)
	if err != nil {
		return nil, err
	}

	// Create owner membership for the current user
	_, err = s.membershipRepo.Create(ctx, userID, createdTenant.ID, model.MembershipRoleOwner, nil)
	if err != nil {
		// TODO: Consider rolling back tenant creation on membership failure
		return nil, err
	}

	// Update member count
	createdTenant.MemberCount = 1

	return createdTenant, nil
}

// UpdateTenant updates a tenant. Requires ADMIN+ role.
func (s *TenantService) UpdateTenant(ctx context.Context, id string, input model.UpdateTenantInput) (*model.Tenant, error) {
	// Check authorization
	_, err := s.requireRole(ctx, id, model.MembershipRoleAdmin)
	if err != nil {
		return nil, err
	}

	return s.tenantRepo.Update(ctx, id, input)
}

// DeleteTenant soft-deletes a tenant. Requires OWNER role.
func (s *TenantService) DeleteTenant(ctx context.Context, id string) (bool, error) {
	// Check authorization
	_, err := s.requireRole(ctx, id, model.MembershipRoleOwner)
	if err != nil {
		return false, err
	}

	err = s.tenantRepo.Delete(ctx, id)
	if err != nil {
		return false, err
	}

	return true, nil
}

// GetTenantMembers retrieves all members of a tenant.
func (s *TenantService) GetTenantMembers(ctx context.Context, tenantID string) ([]*model.Membership, error) {
	// Optional: Check if user is a member of the tenant
	// For now, allow anyone to view members
	return s.membershipRepo.FindByTenantID(ctx, tenantID)
}

// InviteMember invites a user to a tenant. Requires ADMIN+ role.
func (s *TenantService) InviteMember(ctx context.Context, tenantID string, input model.InviteMemberInput) (*model.Membership, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	// Check authorization
	_, err = s.requireRole(ctx, tenantID, model.MembershipRoleAdmin)
	if err != nil {
		return nil, err
	}

	// Find the user to invite
	invitee, err := s.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, errors.ErrUserNotFound
	}

	// Set default role if not provided
	role := model.MembershipRoleMember
	if input.Role != nil {
		role = *input.Role
	}

	// Admins cannot invite owners
	inviterMembership, _ := s.membershipRepo.FindByUserAndTenant(ctx, userID, tenantID)
	if role == model.MembershipRoleOwner && inviterMembership.Role != model.MembershipRoleOwner {
		return nil, errors.ErrForbidden
	}

	// Create membership
	return s.membershipRepo.Create(ctx, invitee.ID, tenantID, role, &userID)
}

// UpdateMemberRole updates a member's role. Requires OWNER role.
func (s *TenantService) UpdateMemberRole(ctx context.Context, membershipID string, role model.MembershipRole) (*model.Membership, error) {
	// Get the membership to find the tenant
	membership, err := s.membershipRepo.FindByID(ctx, membershipID)
	if err != nil {
		return nil, err
	}

	tenantID := membership.Tenant.ID

	// Check authorization - only owners can change roles
	_, err = s.requireRole(ctx, tenantID, model.MembershipRoleOwner)
	if err != nil {
		return nil, err
	}

	// Cannot demote the last owner
	if membership.Role == model.MembershipRoleOwner && role != model.MembershipRoleOwner {
		ownerCount, err := s.membershipRepo.CountOwners(ctx, tenantID)
		if err != nil {
			return nil, err
		}
		if ownerCount <= 1 {
			return nil, errors.ErrLastOwner
		}
	}

	return s.membershipRepo.UpdateRole(ctx, membershipID, role)
}

// RemoveMember removes a member from a tenant. Requires ADMIN+ role.
func (s *TenantService) RemoveMember(ctx context.Context, membershipID string) (bool, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return false, err
	}

	// Get the membership to find the tenant and check constraints
	membership, err := s.membershipRepo.FindByID(ctx, membershipID)
	if err != nil {
		return false, err
	}

	tenantID := membership.Tenant.ID

	// Get current user's membership
	currentMembership, err := s.requireRole(ctx, tenantID, model.MembershipRoleAdmin)
	if err != nil {
		return false, err
	}

	// Cannot remove yourself (use LeaveTenant instead)
	if membership.User.ID == userID {
		return false, errors.NewValidationError("membership", "use leaveTenant to remove yourself")
	}

	// Admins cannot remove other admins or owners
	if currentMembership.Role == model.MembershipRoleAdmin {
		if membership.Role == model.MembershipRoleAdmin || membership.Role == model.MembershipRoleOwner {
			return false, errors.ErrForbidden
		}
	}

	// Cannot remove the last owner
	if membership.Role == model.MembershipRoleOwner {
		ownerCount, err := s.membershipRepo.CountOwners(ctx, tenantID)
		if err != nil {
			return false, err
		}
		if ownerCount <= 1 {
			return false, errors.ErrLastOwner
		}
	}

	err = s.membershipRepo.Delete(ctx, membershipID)
	if err != nil {
		return false, err
	}

	return true, nil
}

// LeaveTenant removes the current user from a tenant.
func (s *TenantService) LeaveTenant(ctx context.Context, tenantID string) (bool, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return false, err
	}

	// Get the user's membership
	membership, err := s.membershipRepo.FindByUserAndTenant(ctx, userID, tenantID)
	if err != nil {
		return false, err
	}

	// Cannot leave if you're the last owner
	if membership.Role == model.MembershipRoleOwner {
		ownerCount, err := s.membershipRepo.CountOwners(ctx, tenantID)
		if err != nil {
			return false, err
		}
		if ownerCount <= 1 {
			return false, errors.ErrCannotLeave
		}
	}

	err = s.membershipRepo.Delete(ctx, membership.ID)
	if err != nil {
		return false, err
	}

	return true, nil
}

// Ensure TenantService implements ITenantService
var _ ITenantService = (*TenantService)(nil)
