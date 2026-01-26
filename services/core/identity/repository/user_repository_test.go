package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourusername/grgn-stack/pkg/errors"
	"github.com/yourusername/grgn-stack/services/core/shared/generated/graphql/model"
)

func TestMockUserRepository_FindByID_Success(t *testing.T) {
	// Arrange
	repo := NewMockUserRepository()
	user := &model.User{
		ID:        "user-123",
		Email:     "test@example.com",
		Status:    model.UserStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repo.AddUser(user)

	// Act
	result, err := repo.FindByID(context.Background(), "user-123")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "user-123", result.ID)
	assert.Equal(t, "test@example.com", result.Email)
}

func TestMockUserRepository_FindByID_NotFound(t *testing.T) {
	// Arrange
	repo := NewMockUserRepository()

	// Act
	result, err := repo.FindByID(context.Background(), "nonexistent")

	// Assert
	assert.Nil(t, result)
	assert.ErrorIs(t, err, errors.ErrUserNotFound)
}

func TestMockUserRepository_FindByID_DeletedUser(t *testing.T) {
	// Arrange
	repo := NewMockUserRepository()
	user := &model.User{
		ID:     "user-123",
		Email:  "test@example.com",
		Status: model.UserStatusDeleted,
	}
	repo.AddUser(user)

	// Act
	result, err := repo.FindByID(context.Background(), "user-123")

	// Assert
	assert.Nil(t, result)
	assert.ErrorIs(t, err, errors.ErrUserNotFound)
}

func TestMockUserRepository_FindByEmail_Success(t *testing.T) {
	// Arrange
	repo := NewMockUserRepository()
	user := &model.User{
		ID:     "user-123",
		Email:  "test@example.com",
		Status: model.UserStatusActive,
	}
	repo.AddUser(user)

	// Act
	result, err := repo.FindByEmail(context.Background(), "test@example.com")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "user-123", result.ID)
}

func TestMockUserRepository_FindByEmail_NotFound(t *testing.T) {
	// Arrange
	repo := NewMockUserRepository()

	// Act
	result, err := repo.FindByEmail(context.Background(), "nonexistent@example.com")

	// Assert
	assert.Nil(t, result)
	assert.ErrorIs(t, err, errors.ErrUserNotFound)
}

func TestMockUserRepository_Create_Success(t *testing.T) {
	// Arrange
	repo := NewMockUserRepository()
	name := "Test User"
	user := &model.User{
		Email: "test@example.com",
		Name:  &name,
	}

	// Act
	result, err := repo.Create(context.Background(), user)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, result.ID)
	assert.Equal(t, "test@example.com", result.Email)
	assert.Equal(t, "Test User", *result.Name)
	assert.Equal(t, model.UserStatusActive, result.Status)
	assert.False(t, result.CreatedAt.IsZero())
	assert.False(t, result.UpdatedAt.IsZero())
}

func TestMockUserRepository_Create_DuplicateEmail(t *testing.T) {
	// Arrange
	repo := NewMockUserRepository()
	existingUser := &model.User{
		ID:     "existing-user",
		Email:  "test@example.com",
		Status: model.UserStatusActive,
	}
	repo.AddUser(existingUser)

	newUser := &model.User{
		Email: "test@example.com",
	}

	// Act
	result, err := repo.Create(context.Background(), newUser)

	// Assert
	assert.Nil(t, result)
	assert.ErrorIs(t, err, errors.ErrEmailTaken)
}

func TestMockUserRepository_Update_Success(t *testing.T) {
	// Arrange
	repo := NewMockUserRepository()
	originalName := "Original Name"
	user := &model.User{
		ID:     "user-123",
		Email:  "test@example.com",
		Name:   &originalName,
		Status: model.UserStatusActive,
	}
	repo.AddUser(user)

	newName := "Updated Name"
	input := model.UpdateProfileInput{Name: &newName}

	// Act
	result, err := repo.Update(context.Background(), "user-123", input)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", *result.Name)
}

func TestMockUserRepository_Update_NotFound(t *testing.T) {
	// Arrange
	repo := NewMockUserRepository()
	newName := "Updated Name"
	input := model.UpdateProfileInput{Name: &newName}

	// Act
	result, err := repo.Update(context.Background(), "nonexistent", input)

	// Assert
	assert.Nil(t, result)
	assert.ErrorIs(t, err, errors.ErrUserNotFound)
}

func TestMockUserRepository_Delete_Success(t *testing.T) {
	// Arrange
	repo := NewMockUserRepository()
	user := &model.User{
		ID:     "user-123",
		Email:  "test@example.com",
		Status: model.UserStatusActive,
	}
	repo.AddUser(user)

	// Act
	err := repo.Delete(context.Background(), "user-123")

	// Assert
	require.NoError(t, err)

	// Verify user is now deleted
	_, findErr := repo.FindByID(context.Background(), "user-123")
	assert.ErrorIs(t, findErr, errors.ErrUserNotFound)
}

func TestMockUserRepository_Delete_NotFound(t *testing.T) {
	// Arrange
	repo := NewMockUserRepository()

	// Act
	err := repo.Delete(context.Background(), "nonexistent")

	// Assert
	assert.ErrorIs(t, err, errors.ErrUserNotFound)
}

func TestMockUserRepository_ExistsByEmail(t *testing.T) {
	// Arrange
	repo := NewMockUserRepository()
	user := &model.User{
		ID:     "user-123",
		Email:  "test@example.com",
		Status: model.UserStatusActive,
	}
	repo.AddUser(user)

	// Act & Assert - existing email
	exists, err := repo.ExistsByEmail(context.Background(), "test@example.com")
	require.NoError(t, err)
	assert.True(t, exists)

	// Act & Assert - non-existing email
	exists, err = repo.ExistsByEmail(context.Background(), "nonexistent@example.com")
	require.NoError(t, err)
	assert.False(t, exists)
}

func TestMockUserRepository_List(t *testing.T) {
	// Arrange
	repo := NewMockUserRepository()
	for i := 0; i < 5; i++ {
		repo.AddUser(&model.User{
			ID:     string(rune('a' + i)),
			Email:  string(rune('a'+i)) + "@example.com",
			Status: model.UserStatusActive,
		})
	}

	// Act
	users, err := repo.List(context.Background(), 3, 0)

	// Assert
	require.NoError(t, err)
	assert.Len(t, users, 3)
}

func TestMockUserRepository_List_WithOffset(t *testing.T) {
	// Arrange
	repo := NewMockUserRepository()
	for i := 0; i < 5; i++ {
		repo.AddUser(&model.User{
			ID:     string(rune('a' + i)),
			Email:  string(rune('a'+i)) + "@example.com",
			Status: model.UserStatusActive,
		})
	}

	// Act
	users, err := repo.List(context.Background(), 10, 3)

	// Assert
	require.NoError(t, err)
	assert.Len(t, users, 2) // 5 total - 3 offset = 2 remaining
}
