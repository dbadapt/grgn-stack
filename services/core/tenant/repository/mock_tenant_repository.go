package repository

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/grgn-stack/pkg/errors"
	"github.com/yourusername/grgn-stack/services/core/shared/generated/graphql/model"
)

// MockTenantRepository is a mock implementation of ITenantRepository for testing.
type MockTenantRepository struct {
	mu      sync.RWMutex
	tenants map[string]*model.Tenant

	// Function overrides for testing specific behaviors
	FindByIDFunc       func(ctx context.Context, id string) (*model.Tenant, error)
	FindBySlugFunc     func(ctx context.Context, slug string) (*model.Tenant, error)
	FindByUserIDFunc   func(ctx context.Context, userID string) ([]*model.Tenant, error)
	CreateFunc         func(ctx context.Context, tenant *model.Tenant) (*model.Tenant, error)
	UpdateFunc         func(ctx context.Context, id string, input model.UpdateTenantInput) (*model.Tenant, error)
	DeleteFunc         func(ctx context.Context, id string) error
	ExistsBySlugFunc   func(ctx context.Context, slug string) (bool, error)
	GetMemberCountFunc func(ctx context.Context, tenantID string) (int, error)

	// For testing: track user-tenant relationships
	userTenants map[string][]string // userID -> []tenantID
}

// NewMockTenantRepository creates a new MockTenantRepository.
func NewMockTenantRepository() *MockTenantRepository {
	return &MockTenantRepository{
		tenants:     make(map[string]*model.Tenant),
		userTenants: make(map[string][]string),
	}
}

// AddTenant adds a tenant to the mock repository for testing.
func (m *MockTenantRepository) AddTenant(tenant *model.Tenant) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.tenants[tenant.ID] = tenant
}

// AddUserToTenant associates a user with a tenant for testing.
func (m *MockTenantRepository) AddUserToTenant(userID, tenantID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.userTenants[userID] = append(m.userTenants[userID], tenantID)
}

// Reset clears all data from the mock repository.
func (m *MockTenantRepository) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.tenants = make(map[string]*model.Tenant)
	m.userTenants = make(map[string][]string)
}

// FindByID retrieves a tenant by ID.
func (m *MockTenantRepository) FindByID(ctx context.Context, id string) (*model.Tenant, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	tenant, ok := m.tenants[id]
	if !ok || tenant.Status == model.TenantStatusDeleted {
		return nil, errors.ErrTenantNotFound
	}
	return tenant, nil
}

// FindBySlug retrieves a tenant by slug.
func (m *MockTenantRepository) FindBySlug(ctx context.Context, slug string) (*model.Tenant, error) {
	if m.FindBySlugFunc != nil {
		return m.FindBySlugFunc(ctx, slug)
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, tenant := range m.tenants {
		if tenant.Slug == slug && tenant.Status != model.TenantStatusDeleted {
			return tenant, nil
		}
	}
	return nil, errors.ErrTenantNotFound
}

// FindByUserID retrieves all tenants a user is a member of.
func (m *MockTenantRepository) FindByUserID(ctx context.Context, userID string) ([]*model.Tenant, error) {
	if m.FindByUserIDFunc != nil {
		return m.FindByUserIDFunc(ctx, userID)
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	tenantIDs, ok := m.userTenants[userID]
	if !ok {
		return []*model.Tenant{}, nil
	}

	var tenants []*model.Tenant
	for _, tenantID := range tenantIDs {
		if tenant, ok := m.tenants[tenantID]; ok && tenant.Status != model.TenantStatusDeleted {
			tenants = append(tenants, tenant)
		}
	}
	return tenants, nil
}

// Create creates a new tenant.
func (m *MockTenantRepository) Create(ctx context.Context, tenant *model.Tenant) (*model.Tenant, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, tenant)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Check for duplicate slug
	for _, existing := range m.tenants {
		if existing.Slug == tenant.Slug && existing.Status != model.TenantStatusDeleted {
			return nil, errors.ErrSlugTaken
		}
	}

	// Generate ID if not provided
	if tenant.ID == "" {
		tenant.ID = uuid.New().String()
	}

	now := time.Now()
	tenant.CreatedAt = now
	tenant.UpdatedAt = now

	if tenant.Status == "" {
		tenant.Status = model.TenantStatusActive
	}
	if tenant.Plan == "" {
		tenant.Plan = model.TenantPlanFree
	}
	if tenant.IsolationMode == "" {
		tenant.IsolationMode = model.TenantIsolationModeShared
	}

	m.tenants[tenant.ID] = tenant
	return tenant, nil
}

// Update updates a tenant.
func (m *MockTenantRepository) Update(ctx context.Context, id string, input model.UpdateTenantInput) (*model.Tenant, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, id, input)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	tenant, ok := m.tenants[id]
	if !ok || tenant.Status == model.TenantStatusDeleted {
		return nil, errors.ErrTenantNotFound
	}

	if input.Name != nil {
		tenant.Name = *input.Name
	}
	if input.Plan != nil {
		tenant.Plan = *input.Plan
	}
	if input.Status != nil {
		tenant.Status = *input.Status
	}
	tenant.UpdatedAt = time.Now()

	return tenant, nil
}

// Delete soft-deletes a tenant.
func (m *MockTenantRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	tenant, ok := m.tenants[id]
	if !ok || tenant.Status == model.TenantStatusDeleted {
		return errors.ErrTenantNotFound
	}

	tenant.Status = model.TenantStatusDeleted
	tenant.UpdatedAt = time.Now()
	return nil
}

// ExistsBySlug checks if a tenant with the given slug exists.
func (m *MockTenantRepository) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
	if m.ExistsBySlugFunc != nil {
		return m.ExistsBySlugFunc(ctx, slug)
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, tenant := range m.tenants {
		if tenant.Slug == slug && tenant.Status != model.TenantStatusDeleted {
			return true, nil
		}
	}
	return false, nil
}

// GetMemberCount returns the number of members in a tenant.
func (m *MockTenantRepository) GetMemberCount(ctx context.Context, tenantID string) (int, error) {
	if m.GetMemberCountFunc != nil {
		return m.GetMemberCountFunc(ctx, tenantID)
	}

	// Default: return 0 (tests should override if needed)
	return 0, nil
}

// Ensure MockTenantRepository implements ITenantRepository
var _ ITenantRepository = (*MockTenantRepository)(nil)
