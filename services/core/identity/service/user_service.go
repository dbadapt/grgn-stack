package service

import (
	"context"

	"github.com/yourusername/grgn-stack/pkg/auth"
	"github.com/yourusername/grgn-stack/services/core/identity/repository"
	"github.com/yourusername/grgn-stack/services/core/shared/generated/graphql/model"
)

// UserService implements IUserService with business logic.
type UserService struct {
	userRepo repository.IUserRepository
}

// NewUserService creates a new UserService.
func NewUserService(userRepo repository.IUserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetCurrentUser retrieves the currently authenticated user.
func (s *UserService) GetCurrentUser(ctx context.Context) (*model.User, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	return s.userRepo.FindByID(ctx, userID)
}

// GetUserByID retrieves a user by their ID.
func (s *UserService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

// UpdateProfile updates the current user's profile.
func (s *UserService) UpdateProfile(ctx context.Context, input model.UpdateProfileInput) (*model.User, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	return s.userRepo.Update(ctx, userID, input)
}

// DeleteAccount soft-deletes the current user's account.
func (s *UserService) DeleteAccount(ctx context.Context) error {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return err
	}

	return s.userRepo.Delete(ctx, userID)
}

// CreateUser creates a new user (internal use).
func (s *UserService) CreateUser(ctx context.Context, email string, name *string) (*model.User, error) {
	user := &model.User{
		Email:  email,
		Name:   name,
		Status: model.UserStatusActive,
	}

	return s.userRepo.Create(ctx, user)
}

// GetUserByEmail retrieves a user by email (internal use).
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.userRepo.FindByEmail(ctx, email)
}

// Ensure UserService implements IUserService
var _ IUserService = (*UserService)(nil)
