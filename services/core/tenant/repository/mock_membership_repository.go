package repository

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/grgn-stack/pkg/errors"
	"github.com/yourusername/grgn-stack/services/core/shared/generated/graphql/model"
)

// MockMembershipRepository is a mock implementation of IMembershipRepository for testing.
type MockMembershipRepository struct {
	mu          sync.RWMutex
	memberships map[string]*model.Membership

	// Index maps for efficient lookups
	byTenant map[string][]string // tenantID -> []membershipID
	byUser   map[string][]string // userID -> []membershipID

	// Function overrides for testing specific behaviors
	FindByIDFunc                  func(ctx context.Context, id string) (*model.Membership, error)
	FindByTenantIDFunc            func(ctx context.Context, tenantID string) ([]*model.Membership, error)
	FindByUserIDFunc              func(ctx context.Context, userID string) ([]*model.Membership, error)
	FindByUserAndTenantFunc       func(ctx context.Context, userID, tenantID string) (*model.Membership, error)
	CreateFunc                    func(ctx context.Context, userID, tenantID string, role model.MembershipRole, invitedByID *string) (*model.Membership, error)
	UpdateRoleFunc                func(ctx context.Context, id string, role model.MembershipRole) (*model.Membership, error)
	DeleteFunc                    func(ctx context.Context, id string) error
	CountOwnersFunc               func(ctx context.Context, tenantID string) (int, error)
	GetTenantIDByMembershipIDFunc func(ctx context.Context, membershipID string) (string, error)
	GetUserIDByMembershipIDFunc   func(ctx context.Context, membershipID string) (string, error)
}

// NewMockMembershipRepository creates a new MockMembershipRepository.
func NewMockMembershipRepository() *MockMembershipRepository {
	return &MockMembershipRepository{
		memberships: make(map[string]*model.Membership),
		byTenant:    make(map[string][]string),
		byUser:      make(map[string][]string),
	}
}

// AddMembership adds a membership to the mock repository for testing.
func (m *MockMembershipRepository) AddMembership(membership *model.Membership) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.memberships[membership.ID] = membership

	// Update indexes
	if membership.Tenant != nil {
		m.byTenant[membership.Tenant.ID] = append(m.byTenant[membership.Tenant.ID], membership.ID)
	}
	if membership.User != nil {
		m.byUser[membership.User.ID] = append(m.byUser[membership.User.ID], membership.ID)
	}
}

// Reset clears all data from the mock repository.
func (m *MockMembershipRepository) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.memberships = make(map[string]*model.Membership)
	m.byTenant = make(map[string][]string)
	m.byUser = make(map[string][]string)
}

// FindByID retrieves a membership by ID.
func (m *MockMembershipRepository) FindByID(ctx context.Context, id string) (*model.Membership, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	membership, ok := m.memberships[id]
	if !ok {
		return nil, errors.ErrMembershipNotFound
	}
	return membership, nil
}

// FindByTenantID retrieves all memberships for a tenant.
func (m *MockMembershipRepository) FindByTenantID(ctx context.Context, tenantID string) ([]*model.Membership, error) {
	if m.FindByTenantIDFunc != nil {
		return m.FindByTenantIDFunc(ctx, tenantID)
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	membershipIDs, ok := m.byTenant[tenantID]
	if !ok {
		return []*model.Membership{}, nil
	}

	var memberships []*model.Membership
	for _, id := range membershipIDs {
		if membership, ok := m.memberships[id]; ok {
			memberships = append(memberships, membership)
		}
	}
	return memberships, nil
}

// FindByUserID retrieves all memberships for a user.
func (m *MockMembershipRepository) FindByUserID(ctx context.Context, userID string) ([]*model.Membership, error) {
	if m.FindByUserIDFunc != nil {
		return m.FindByUserIDFunc(ctx, userID)
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	membershipIDs, ok := m.byUser[userID]
	if !ok {
		return []*model.Membership{}, nil
	}

	var memberships []*model.Membership
	for _, id := range membershipIDs {
		if membership, ok := m.memberships[id]; ok {
			memberships = append(memberships, membership)
		}
	}
	return memberships, nil
}

// FindByUserAndTenant retrieves a membership by user and tenant.
func (m *MockMembershipRepository) FindByUserAndTenant(ctx context.Context, userID, tenantID string) (*model.Membership, error) {
	if m.FindByUserAndTenantFunc != nil {
		return m.FindByUserAndTenantFunc(ctx, userID, tenantID)
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, membership := range m.memberships {
		if membership.User != nil && membership.Tenant != nil {
			if membership.User.ID == userID && membership.Tenant.ID == tenantID {
				return membership, nil
			}
		}
	}
	return nil, errors.ErrMembershipNotFound
}

// Create creates a new membership.
func (m *MockMembershipRepository) Create(ctx context.Context, userID, tenantID string, role model.MembershipRole, invitedByID *string) (*model.Membership, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, userID, tenantID, role, invitedByID)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if already a member
	for _, membership := range m.memberships {
		if membership.User != nil && membership.Tenant != nil {
			if membership.User.ID == userID && membership.Tenant.ID == tenantID {
				return nil, errors.ErrAlreadyMember
			}
		}
	}

	membership := &model.Membership{
		ID:       uuid.New().String(),
		Role:     role,
		JoinedAt: time.Now(),
		User:     &model.User{ID: userID},
		Tenant:   &model.Tenant{ID: tenantID},
	}

	if invitedByID != nil && *invitedByID != "" {
		membership.InvitedBy = &model.User{ID: *invitedByID}
	}

	m.memberships[membership.ID] = membership
	m.byTenant[tenantID] = append(m.byTenant[tenantID], membership.ID)
	m.byUser[userID] = append(m.byUser[userID], membership.ID)

	return membership, nil
}

// UpdateRole updates a membership's role.
func (m *MockMembershipRepository) UpdateRole(ctx context.Context, id string, role model.MembershipRole) (*model.Membership, error) {
	if m.UpdateRoleFunc != nil {
		return m.UpdateRoleFunc(ctx, id, role)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	membership, ok := m.memberships[id]
	if !ok {
		return nil, errors.ErrMembershipNotFound
	}

	membership.Role = role
	return membership, nil
}

// Delete removes a membership.
func (m *MockMembershipRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	membership, ok := m.memberships[id]
	if !ok {
		return errors.ErrMembershipNotFound
	}

	// Remove from indexes
	if membership.Tenant != nil {
		m.byTenant[membership.Tenant.ID] = m.removeFromSlice(m.byTenant[membership.Tenant.ID], id)
	}
	if membership.User != nil {
		m.byUser[membership.User.ID] = m.removeFromSlice(m.byUser[membership.User.ID], id)
	}

	delete(m.memberships, id)
	return nil
}

// CountOwners returns the number of owners in a tenant.
func (m *MockMembershipRepository) CountOwners(ctx context.Context, tenantID string) (int, error) {
	if m.CountOwnersFunc != nil {
		return m.CountOwnersFunc(ctx, tenantID)
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	count := 0
	membershipIDs, ok := m.byTenant[tenantID]
	if !ok {
		return 0, nil
	}

	for _, id := range membershipIDs {
		if membership, ok := m.memberships[id]; ok {
			if membership.Role == model.MembershipRoleOwner {
				count++
			}
		}
	}
	return count, nil
}

// GetTenantIDByMembershipID returns the tenant ID for a membership.
func (m *MockMembershipRepository) GetTenantIDByMembershipID(ctx context.Context, membershipID string) (string, error) {
	if m.GetTenantIDByMembershipIDFunc != nil {
		return m.GetTenantIDByMembershipIDFunc(ctx, membershipID)
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	membership, ok := m.memberships[membershipID]
	if !ok || membership.Tenant == nil {
		return "", errors.ErrMembershipNotFound
	}
	return membership.Tenant.ID, nil
}

// GetUserIDByMembershipID returns the user ID for a membership.
func (m *MockMembershipRepository) GetUserIDByMembershipID(ctx context.Context, membershipID string) (string, error) {
	if m.GetUserIDByMembershipIDFunc != nil {
		return m.GetUserIDByMembershipIDFunc(ctx, membershipID)
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	membership, ok := m.memberships[membershipID]
	if !ok || membership.User == nil {
		return "", errors.ErrMembershipNotFound
	}
	return membership.User.ID, nil
}

// removeFromSlice removes an element from a slice and returns the new slice.
func (m *MockMembershipRepository) removeFromSlice(slice []string, item string) []string {
	for i, v := range slice {
		if v == item {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

// Ensure MockMembershipRepository implements IMembershipRepository
var _ IMembershipRepository = (*MockMembershipRepository)(nil)
