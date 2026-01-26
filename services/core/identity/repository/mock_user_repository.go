package repository

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/grgn-stack/pkg/errors"
	"github.com/yourusername/grgn-stack/services/core/shared/generated/graphql/model"
)

// MockUserRepository is a mock implementation of IUserRepository for testing.
type MockUserRepository struct {
	mu    sync.RWMutex
	users map[string]*model.User

	// Function overrides for testing specific behaviors
	FindByIDFunc      func(ctx context.Context, id string) (*model.User, error)
	FindByEmailFunc   func(ctx context.Context, email string) (*model.User, error)
	CreateFunc        func(ctx context.Context, user *model.User) (*model.User, error)
	UpdateFunc        func(ctx context.Context, id string, input model.UpdateProfileInput) (*model.User, error)
	DeleteFunc        func(ctx context.Context, id string) error
	ListFunc          func(ctx context.Context, limit, offset int) ([]*model.User, error)
	ExistsByEmailFunc func(ctx context.Context, email string) (bool, error)
}

// NewMockUserRepository creates a new MockUserRepository.
func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]*model.User),
	}
}

// AddUser adds a user to the mock repository for testing.
func (m *MockUserRepository) AddUser(user *model.User) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.users[user.ID] = user
}

// GetUsers returns all users in the mock repository.
func (m *MockUserRepository) GetUsers() map[string]*model.User {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.users
}

// Reset clears all users from the mock repository.
func (m *MockUserRepository) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.users = make(map[string]*model.User)
}

// FindByID retrieves a user by ID.
func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	user, ok := m.users[id]
	if !ok || user.Status == model.UserStatusDeleted {
		return nil, errors.ErrUserNotFound
	}
	return user, nil
}

// FindByEmail retrieves a user by email.
func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	if m.FindByEmailFunc != nil {
		return m.FindByEmailFunc(ctx, email)
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, user := range m.users {
		if user.Email == email && user.Status != model.UserStatusDeleted {
			return user, nil
		}
	}
	return nil, errors.ErrUserNotFound
}

// Create creates a new user.
func (m *MockUserRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, user)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Check for duplicate email
	for _, existing := range m.users {
		if existing.Email == user.Email && existing.Status != model.UserStatusDeleted {
			return nil, errors.ErrEmailTaken
		}
	}

	// Generate ID if not provided
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	if user.Status == "" {
		user.Status = model.UserStatusActive
	}

	m.users[user.ID] = user
	return user, nil
}

// Update updates a user's profile.
func (m *MockUserRepository) Update(ctx context.Context, id string, input model.UpdateProfileInput) (*model.User, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, id, input)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	user, ok := m.users[id]
	if !ok || user.Status == model.UserStatusDeleted {
		return nil, errors.ErrUserNotFound
	}

	if input.Name != nil {
		user.Name = input.Name
	}
	if input.AvatarURL != nil {
		user.AvatarURL = input.AvatarURL
	}
	user.UpdatedAt = time.Now()

	return user, nil
}

// Delete soft-deletes a user.
func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	user, ok := m.users[id]
	if !ok || user.Status == model.UserStatusDeleted {
		return errors.ErrUserNotFound
	}

	user.Status = model.UserStatusDeleted
	user.UpdatedAt = time.Now()
	return nil
}

// List retrieves users with pagination.
func (m *MockUserRepository) List(ctx context.Context, limit, offset int) ([]*model.User, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, limit, offset)
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	var users []*model.User
	for _, user := range m.users {
		if user.Status != model.UserStatusDeleted {
			users = append(users, user)
		}
	}

	// Apply pagination
	start := offset
	if start > len(users) {
		return []*model.User{}, nil
	}

	end := start + limit
	if end > len(users) {
		end = len(users)
	}

	return users[start:end], nil
}

// ExistsByEmail checks if a user with the given email exists.
func (m *MockUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	if m.ExistsByEmailFunc != nil {
		return m.ExistsByEmailFunc(ctx, email)
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, user := range m.users {
		if user.Email == email && user.Status != model.UserStatusDeleted {
			return true, nil
		}
	}
	return false, nil
}

// Ensure MockUserRepository implements IUserRepository
var _ IUserRepository = (*MockUserRepository)(nil)
