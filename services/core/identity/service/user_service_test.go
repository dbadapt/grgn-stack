package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourusername/grgn-stack/pkg/auth"
	"github.com/yourusername/grgn-stack/pkg/errors"
	"github.com/yourusername/grgn-stack/services/core/identity/repository"
	"github.com/yourusername/grgn-stack/services/core/shared/generated/graphql/model"
)

func TestUserService_GetCurrentUser_Success(t *testing.T) {
	// Arrange
	mockRepo := repository.NewMockUserRepository()
	mockRepo.AddUser(&model.User{
		ID:        "user-123",
		Email:     "test@example.com",
		Status:    model.UserStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	svc := NewUserService(mockRepo)
	ctx := auth.WithUserID(context.Background(), "user-123")

	// Act
	user, err := svc.GetCurrentUser(ctx)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "user-123", user.ID)
	assert.Equal(t, "test@example.com", user.Email)
}

func TestUserService_GetCurrentUser_NotAuthenticated(t *testing.T) {
	// Arrange
	mockRepo := repository.NewMockUserRepository()
	svc := NewUserService(mockRepo)
	ctx := context.Background() // No user in context

	// Act
	user, err := svc.GetCurrentUser(ctx)

	// Assert
	assert.Nil(t, user)
	assert.ErrorIs(t, err, errors.ErrNotAuthenticated)
}

func TestUserService_GetCurrentUser_UserNotFound(t *testing.T) {
	// Arrange
	mockRepo := repository.NewMockUserRepository()
	// No user added to mock

	svc := NewUserService(mockRepo)
	ctx := auth.WithUserID(context.Background(), "nonexistent")

	// Act
	user, err := svc.GetCurrentUser(ctx)

	// Assert
	assert.Nil(t, user)
	assert.ErrorIs(t, err, errors.ErrUserNotFound)
}

func TestUserService_GetUserByID_Success(t *testing.T) {
	// Arrange
	mockRepo := repository.NewMockUserRepository()
	mockRepo.AddUser(&model.User{
		ID:     "user-123",
		Email:  "test@example.com",
		Status: model.UserStatusActive,
	})

	svc := NewUserService(mockRepo)

	// Act
	user, err := svc.GetUserByID(context.Background(), "user-123")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "user-123", user.ID)
}

func TestUserService_GetUserByID_NotFound(t *testing.T) {
	// Arrange
	mockRepo := repository.NewMockUserRepository()
	svc := NewUserService(mockRepo)

	// Act
	user, err := svc.GetUserByID(context.Background(), "nonexistent")

	// Assert
	assert.Nil(t, user)
	assert.ErrorIs(t, err, errors.ErrUserNotFound)
}

func TestUserService_UpdateProfile_Success(t *testing.T) {
	// Arrange
	originalName := "Original Name"
	mockRepo := repository.NewMockUserRepository()
	mockRepo.AddUser(&model.User{
		ID:     "user-123",
		Email:  "test@example.com",
		Name:   &originalName,
		Status: model.UserStatusActive,
	})

	svc := NewUserService(mockRepo)
	ctx := auth.WithUserID(context.Background(), "user-123")

	newName := "Updated Name"
	input := model.UpdateProfileInput{Name: &newName}

	// Act
	user, err := svc.UpdateProfile(ctx, input)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", *user.Name)
}

func TestUserService_UpdateProfile_NotAuthenticated(t *testing.T) {
	// Arrange
	mockRepo := repository.NewMockUserRepository()
	svc := NewUserService(mockRepo)
	ctx := context.Background()

	newName := "Updated Name"
	input := model.UpdateProfileInput{Name: &newName}

	// Act
	user, err := svc.UpdateProfile(ctx, input)

	// Assert
	assert.Nil(t, user)
	assert.ErrorIs(t, err, errors.ErrNotAuthenticated)
}

func TestUserService_UpdateProfile_UserNotFound(t *testing.T) {
	// Arrange
	mockRepo := repository.NewMockUserRepository()
	svc := NewUserService(mockRepo)
	ctx := auth.WithUserID(context.Background(), "nonexistent")

	newName := "Updated Name"
	input := model.UpdateProfileInput{Name: &newName}

	// Act
	user, err := svc.UpdateProfile(ctx, input)

	// Assert
	assert.Nil(t, user)
	assert.ErrorIs(t, err, errors.ErrUserNotFound)
}

func TestUserService_DeleteAccount_Success(t *testing.T) {
	// Arrange
	mockRepo := repository.NewMockUserRepository()
	mockRepo.AddUser(&model.User{
		ID:     "user-123",
		Email:  "test@example.com",
		Status: model.UserStatusActive,
	})

	svc := NewUserService(mockRepo)
	ctx := auth.WithUserID(context.Background(), "user-123")

	// Act
	err := svc.DeleteAccount(ctx)

	// Assert
	require.NoError(t, err)

	// Verify user is now deleted
	_, findErr := mockRepo.FindByID(context.Background(), "user-123")
	assert.ErrorIs(t, findErr, errors.ErrUserNotFound)
}

func TestUserService_DeleteAccount_NotAuthenticated(t *testing.T) {
	// Arrange
	mockRepo := repository.NewMockUserRepository()
	svc := NewUserService(mockRepo)
	ctx := context.Background()

	// Act
	err := svc.DeleteAccount(ctx)

	// Assert
	assert.ErrorIs(t, err, errors.ErrNotAuthenticated)
}

func TestUserService_CreateUser_Success(t *testing.T) {
	// Arrange
	mockRepo := repository.NewMockUserRepository()
	svc := NewUserService(mockRepo)

	name := "Test User"

	// Act
	user, err := svc.CreateUser(context.Background(), "test@example.com", &name)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, user.ID)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "Test User", *user.Name)
	assert.Equal(t, model.UserStatusActive, user.Status)
}

func TestUserService_CreateUser_DuplicateEmail(t *testing.T) {
	// Arrange
	mockRepo := repository.NewMockUserRepository()
	mockRepo.AddUser(&model.User{
		ID:     "existing-user",
		Email:  "test@example.com",
		Status: model.UserStatusActive,
	})

	svc := NewUserService(mockRepo)
	name := "New User"

	// Act
	user, err := svc.CreateUser(context.Background(), "test@example.com", &name)

	// Assert
	assert.Nil(t, user)
	assert.ErrorIs(t, err, errors.ErrEmailTaken)
}

func TestUserService_GetUserByEmail_Success(t *testing.T) {
	// Arrange
	mockRepo := repository.NewMockUserRepository()
	mockRepo.AddUser(&model.User{
		ID:     "user-123",
		Email:  "test@example.com",
		Status: model.UserStatusActive,
	})

	svc := NewUserService(mockRepo)

	// Act
	user, err := svc.GetUserByEmail(context.Background(), "test@example.com")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "user-123", user.ID)
}

func TestUserService_GetUserByEmail_NotFound(t *testing.T) {
	// Arrange
	mockRepo := repository.NewMockUserRepository()
	svc := NewUserService(mockRepo)

	// Act
	user, err := svc.GetUserByEmail(context.Background(), "nonexistent@example.com")

	// Assert
	assert.Nil(t, user)
	assert.ErrorIs(t, err, errors.ErrUserNotFound)
}
