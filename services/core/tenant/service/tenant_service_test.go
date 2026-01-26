package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourusername/grgn-stack/pkg/auth"
	"github.com/yourusername/grgn-stack/pkg/errors"
	identityRepo "github.com/yourusername/grgn-stack/services/core/identity/repository"
	"github.com/yourusername/grgn-stack/services/core/shared/generated/graphql/model"
	"github.com/yourusername/grgn-stack/services/core/tenant/repository"
)

func setupTestService() (*TenantService, *repository.MockTenantRepository, *repository.MockMembershipRepository, *identityRepo.MockUserRepository) {
	tenantRepo := repository.NewMockTenantRepository()
	membershipRepo := repository.NewMockMembershipRepository()
	userRepo := identityRepo.NewMockUserRepository()

	svc := NewTenantService(tenantRepo, membershipRepo, userRepo)
	return svc, tenantRepo, membershipRepo, userRepo
}

func TestTenantService_CreateTenant_Success(t *testing.T) {
	// Arrange
	svc, _, membershipRepo, _ := setupTestService()
	ctx := auth.WithUserID(context.Background(), "user-123")

	input := model.CreateTenantInput{
		Name: "Acme Corp",
		Slug: "acme-corp",
	}

	// Act
	tenant, err := svc.CreateTenant(ctx, input)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, tenant.ID)
	assert.Equal(t, "Acme Corp", tenant.Name)
	assert.Equal(t, "acme-corp", tenant.Slug)
	assert.Equal(t, model.TenantPlanFree, tenant.Plan)
	assert.Equal(t, model.TenantStatusActive, tenant.Status)
	assert.Equal(t, 1, tenant.MemberCount)

	// Verify owner membership was created
	memberships, _ := membershipRepo.FindByTenantID(ctx, tenant.ID)
	require.Len(t, memberships, 1)
	assert.Equal(t, model.MembershipRoleOwner, memberships[0].Role)
	assert.Equal(t, "user-123", memberships[0].User.ID)
}

func TestTenantService_CreateTenant_NotAuthenticated(t *testing.T) {
	// Arrange
	svc, _, _, _ := setupTestService()
	ctx := context.Background()

	input := model.CreateTenantInput{
		Name: "Acme Corp",
		Slug: "acme",
	}

	// Act
	tenant, err := svc.CreateTenant(ctx, input)

	// Assert
	assert.Nil(t, tenant)
	assert.ErrorIs(t, err, errors.ErrNotAuthenticated)
}

func TestTenantService_CreateTenant_InvalidSlug(t *testing.T) {
	// Arrange
	svc, _, _, _ := setupTestService()
	ctx := auth.WithUserID(context.Background(), "user-123")

	testCases := []struct {
		slug string
		desc string
	}{
		{"ab", "too short"},
		{"has spaces", "contains spaces"},
		{"has.dots", "contains dots"},
		{"special!chars", "contains special chars"},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			input := model.CreateTenantInput{Name: "Test", Slug: tc.slug}
			tenant, err := svc.CreateTenant(ctx, input)
			assert.Nil(t, tenant)
			assert.ErrorIs(t, err, errors.ErrInvalidSlug)
		})
	}
}

func TestTenantService_CreateTenant_ValidSlugs(t *testing.T) {
	// Arrange
	svc, _, _, _ := setupTestService()
	ctx := auth.WithUserID(context.Background(), "user-123")

	validSlugs := []string{
		"abc",
		"ABC",
		"abc-123",
		"abc_123",
		"ABC-123_xyz",
		"a1b2c3",
	}

	for _, slug := range validSlugs {
		t.Run(slug, func(t *testing.T) {
			input := model.CreateTenantInput{Name: "Test", Slug: slug}
			tenant, err := svc.CreateTenant(ctx, input)
			require.NoError(t, err)
			assert.Equal(t, slug, tenant.Slug)
		})
	}
}

func TestTenantService_CreateTenant_DuplicateSlug(t *testing.T) {
	// Arrange
	svc, tenantRepo, _, _ := setupTestService()
	ctx := auth.WithUserID(context.Background(), "user-123")

	// Add existing tenant
	tenantRepo.AddTenant(&model.Tenant{
		ID:     "existing-tenant",
		Slug:   "acme",
		Status: model.TenantStatusActive,
	})

	input := model.CreateTenantInput{
		Name: "Another Acme",
		Slug: "acme",
	}

	// Act
	tenant, err := svc.CreateTenant(ctx, input)

	// Assert
	assert.Nil(t, tenant)
	assert.ErrorIs(t, err, errors.ErrSlugTaken)
}

func TestTenantService_GetMyTenants(t *testing.T) {
	// Arrange
	svc, tenantRepo, membershipRepo, _ := setupTestService()
	ctx := auth.WithUserID(context.Background(), "user-123")

	// Add tenants
	tenant1 := &model.Tenant{ID: "tenant-1", Name: "Tenant 1", Slug: "tenant-1", Status: model.TenantStatusActive}
	tenant2 := &model.Tenant{ID: "tenant-2", Name: "Tenant 2", Slug: "tenant-2", Status: model.TenantStatusActive}
	tenantRepo.AddTenant(tenant1)
	tenantRepo.AddTenant(tenant2)
	tenantRepo.AddUserToTenant("user-123", "tenant-1")
	tenantRepo.AddUserToTenant("user-123", "tenant-2")

	// Add memberships
	membershipRepo.AddMembership(&model.Membership{
		ID:     "m1",
		Role:   model.MembershipRoleOwner,
		User:   &model.User{ID: "user-123"},
		Tenant: tenant1,
	})
	membershipRepo.AddMembership(&model.Membership{
		ID:     "m2",
		Role:   model.MembershipRoleMember,
		User:   &model.User{ID: "user-123"},
		Tenant: tenant2,
	})

	// Act
	tenants, err := svc.GetMyTenants(ctx)

	// Assert
	require.NoError(t, err)
	assert.Len(t, tenants, 2)
}

func TestTenantService_UpdateTenant_Success(t *testing.T) {
	// Arrange
	svc, tenantRepo, membershipRepo, _ := setupTestService()
	ctx := auth.WithUserID(context.Background(), "user-123")

	tenant := &model.Tenant{ID: "tenant-1", Name: "Old Name", Slug: "tenant-1", Status: model.TenantStatusActive}
	tenantRepo.AddTenant(tenant)

	// Add admin membership
	membershipRepo.AddMembership(&model.Membership{
		ID:     "m1",
		Role:   model.MembershipRoleAdmin,
		User:   &model.User{ID: "user-123"},
		Tenant: tenant,
	})

	newName := "New Name"
	input := model.UpdateTenantInput{Name: &newName}

	// Act
	updated, err := svc.UpdateTenant(ctx, "tenant-1", input)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "New Name", updated.Name)
}

func TestTenantService_UpdateTenant_NotAdmin(t *testing.T) {
	// Arrange
	svc, tenantRepo, membershipRepo, _ := setupTestService()
	ctx := auth.WithUserID(context.Background(), "user-123")

	tenant := &model.Tenant{ID: "tenant-1", Name: "Tenant", Slug: "tenant", Status: model.TenantStatusActive}
	tenantRepo.AddTenant(tenant)

	// Add member (not admin) membership
	membershipRepo.AddMembership(&model.Membership{
		ID:     "m1",
		Role:   model.MembershipRoleMember,
		User:   &model.User{ID: "user-123"},
		Tenant: tenant,
	})

	newName := "New Name"
	input := model.UpdateTenantInput{Name: &newName}

	// Act
	updated, err := svc.UpdateTenant(ctx, "tenant-1", input)

	// Assert
	assert.Nil(t, updated)
	assert.ErrorIs(t, err, errors.ErrForbidden)
}

func TestTenantService_DeleteTenant_Success(t *testing.T) {
	// Arrange
	svc, tenantRepo, membershipRepo, _ := setupTestService()
	ctx := auth.WithUserID(context.Background(), "user-123")

	tenant := &model.Tenant{ID: "tenant-1", Name: "Tenant", Slug: "tenant", Status: model.TenantStatusActive}
	tenantRepo.AddTenant(tenant)

	// Add owner membership
	membershipRepo.AddMembership(&model.Membership{
		ID:     "m1",
		Role:   model.MembershipRoleOwner,
		User:   &model.User{ID: "user-123"},
		Tenant: tenant,
	})

	// Act
	deleted, err := svc.DeleteTenant(ctx, "tenant-1")

	// Assert
	require.NoError(t, err)
	assert.True(t, deleted)

	// Verify tenant is deleted
	_, findErr := tenantRepo.FindByID(ctx, "tenant-1")
	assert.ErrorIs(t, findErr, errors.ErrTenantNotFound)
}

func TestTenantService_DeleteTenant_NotOwner(t *testing.T) {
	// Arrange
	svc, tenantRepo, membershipRepo, _ := setupTestService()
	ctx := auth.WithUserID(context.Background(), "user-123")

	tenant := &model.Tenant{ID: "tenant-1", Name: "Tenant", Slug: "tenant", Status: model.TenantStatusActive}
	tenantRepo.AddTenant(tenant)

	// Add admin (not owner) membership
	membershipRepo.AddMembership(&model.Membership{
		ID:     "m1",
		Role:   model.MembershipRoleAdmin,
		User:   &model.User{ID: "user-123"},
		Tenant: tenant,
	})

	// Act
	deleted, err := svc.DeleteTenant(ctx, "tenant-1")

	// Assert
	assert.False(t, deleted)
	assert.ErrorIs(t, err, errors.ErrForbidden)
}

func TestTenantService_InviteMember_Success(t *testing.T) {
	// Arrange
	svc, tenantRepo, membershipRepo, userRepo := setupTestService()
	ctx := auth.WithUserID(context.Background(), "admin-123")

	tenant := &model.Tenant{ID: "tenant-1", Name: "Tenant", Slug: "tenant", Status: model.TenantStatusActive}
	tenantRepo.AddTenant(tenant)

	// Add admin membership
	membershipRepo.AddMembership(&model.Membership{
		ID:     "m1",
		Role:   model.MembershipRoleAdmin,
		User:   &model.User{ID: "admin-123"},
		Tenant: tenant,
	})

	// Add user to invite
	userRepo.AddUser(&model.User{ID: "invitee-123", Email: "invitee@example.com", Status: model.UserStatusActive})

	input := model.InviteMemberInput{Email: "invitee@example.com"}

	// Act
	membership, err := svc.InviteMember(ctx, "tenant-1", input)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, model.MembershipRoleMember, membership.Role)
	assert.Equal(t, "invitee-123", membership.User.ID)
}

func TestTenantService_InviteMember_UserNotFound(t *testing.T) {
	// Arrange
	svc, tenantRepo, membershipRepo, _ := setupTestService()
	ctx := auth.WithUserID(context.Background(), "admin-123")

	tenant := &model.Tenant{ID: "tenant-1", Name: "Tenant", Slug: "tenant", Status: model.TenantStatusActive}
	tenantRepo.AddTenant(tenant)

	membershipRepo.AddMembership(&model.Membership{
		ID:     "m1",
		Role:   model.MembershipRoleAdmin,
		User:   &model.User{ID: "admin-123"},
		Tenant: tenant,
	})

	input := model.InviteMemberInput{Email: "nonexistent@example.com"}

	// Act
	membership, err := svc.InviteMember(ctx, "tenant-1", input)

	// Assert
	assert.Nil(t, membership)
	assert.ErrorIs(t, err, errors.ErrUserNotFound)
}

func TestTenantService_LeaveTenant_Success(t *testing.T) {
	// Arrange
	svc, tenantRepo, membershipRepo, _ := setupTestService()
	ctx := auth.WithUserID(context.Background(), "user-123")

	tenant := &model.Tenant{ID: "tenant-1", Name: "Tenant", Slug: "tenant", Status: model.TenantStatusActive}
	tenantRepo.AddTenant(tenant)

	// Add two owners so one can leave
	membershipRepo.AddMembership(&model.Membership{
		ID:     "m1",
		Role:   model.MembershipRoleOwner,
		User:   &model.User{ID: "user-123"},
		Tenant: tenant,
	})
	membershipRepo.AddMembership(&model.Membership{
		ID:     "m2",
		Role:   model.MembershipRoleOwner,
		User:   &model.User{ID: "other-owner"},
		Tenant: tenant,
	})

	// Act
	left, err := svc.LeaveTenant(ctx, "tenant-1")

	// Assert
	require.NoError(t, err)
	assert.True(t, left)

	// Verify membership is deleted
	_, findErr := membershipRepo.FindByUserAndTenant(ctx, "user-123", "tenant-1")
	assert.ErrorIs(t, findErr, errors.ErrMembershipNotFound)
}

func TestTenantService_LeaveTenant_LastOwner(t *testing.T) {
	// Arrange
	svc, tenantRepo, membershipRepo, _ := setupTestService()
	ctx := auth.WithUserID(context.Background(), "user-123")

	tenant := &model.Tenant{ID: "tenant-1", Name: "Tenant", Slug: "tenant", Status: model.TenantStatusActive}
	tenantRepo.AddTenant(tenant)

	// Add only one owner
	membershipRepo.AddMembership(&model.Membership{
		ID:     "m1",
		Role:   model.MembershipRoleOwner,
		User:   &model.User{ID: "user-123"},
		Tenant: tenant,
	})

	// Act
	left, err := svc.LeaveTenant(ctx, "tenant-1")

	// Assert
	assert.False(t, left)
	assert.ErrorIs(t, err, errors.ErrCannotLeave)
}

func TestTenantService_UpdateMemberRole_CannotDemoteLastOwner(t *testing.T) {
	// Arrange
	svc, tenantRepo, membershipRepo, _ := setupTestService()
	ctx := auth.WithUserID(context.Background(), "owner-123")

	tenant := &model.Tenant{ID: "tenant-1", Name: "Tenant", Slug: "tenant", Status: model.TenantStatusActive}
	tenantRepo.AddTenant(tenant)

	// Add only one owner
	membershipRepo.AddMembership(&model.Membership{
		ID:     "m1",
		Role:   model.MembershipRoleOwner,
		User:   &model.User{ID: "owner-123"},
		Tenant: tenant,
	})

	// Act - try to demote ourselves
	_, err := svc.UpdateMemberRole(ctx, "m1", model.MembershipRoleAdmin)

	// Assert
	assert.ErrorIs(t, err, errors.ErrLastOwner)
}

func TestTenantService_RemoveMember_CannotRemoveLastOwner(t *testing.T) {
	// Arrange
	svc, tenantRepo, membershipRepo, _ := setupTestService()
	ctx := auth.WithUserID(context.Background(), "other-admin")

	tenant := &model.Tenant{ID: "tenant-1", Name: "Tenant", Slug: "tenant", Status: model.TenantStatusActive}
	tenantRepo.AddTenant(tenant)

	// Add owner
	membershipRepo.AddMembership(&model.Membership{
		ID:     "owner-membership",
		Role:   model.MembershipRoleOwner,
		User:   &model.User{ID: "owner-123"},
		Tenant: tenant,
	})

	// Add another owner who is trying to remove
	membershipRepo.AddMembership(&model.Membership{
		ID:     "admin-membership",
		Role:   model.MembershipRoleOwner,
		User:   &model.User{ID: "other-admin"},
		Tenant: tenant,
	})

	// Set the count to 1 for this test (simulating single owner scenario)
	membershipRepo.CountOwnersFunc = func(ctx context.Context, tenantID string) (int, error) {
		return 1, nil
	}

	// Act
	removed, err := svc.RemoveMember(ctx, "owner-membership")

	// Assert
	assert.False(t, removed)
	assert.ErrorIs(t, err, errors.ErrLastOwner)
}
