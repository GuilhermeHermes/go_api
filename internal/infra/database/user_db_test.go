package database

import (
	"fmt"
	"testing"

	"github/GuilhermeHermes/GO_API/internal/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&entity.User{})
	require.NoError(t, err)

	return db
}

func createTestUser(t *testing.T) *entity.User {
	user, err := entity.NewUser("testuser", "test@example.com", "password123", "user")
	require.NoError(t, err)
	return user
}

func createTestUserWithEmail(t *testing.T, email string) *entity.User {
	user, err := entity.NewUser("testuser", email, "password123", "user")
	require.NoError(t, err)
	return user
}

func TestNewUserRepository(t *testing.T) {
	db := setupTestDB(t)
	userRepo := NewUserRepository(db)

	assert.NotNil(t, userRepo)
	assert.Equal(t, db, userRepo.DB)
}

func TestUser_Create(t *testing.T) {
	t.Run("should create user successfully", func(t *testing.T) {
		db := setupTestDB(t)
		userRepo := NewUserRepository(db)
		user := createTestUserWithEmail(t, "create_test@example.com")

		err := userRepo.Create(user)

		assert.NoError(t, err)
		assert.NotEmpty(t, user.ID)
		assert.NotEmpty(t, user.CreatedAt)
		assert.NotEmpty(t, user.UpdatedAt)

		// Verify user was actually saved to database
		var count int64
		db.Model(&entity.User{}).Where("email = ?", "create_test@example.com").Count(&count)
		assert.Equal(t, int64(1), count, "User should be saved in database")

		// Verify we can find the user
		foundUser, err := userRepo.FindByEmail("create_test@example.com")
		assert.NoError(t, err)
		assert.NotNil(t, foundUser)
		assert.Equal(t, user.ID, foundUser.ID)
	})

	t.Run("should not allow duplicate email", func(t *testing.T) {
		db := setupTestDB(t)
		userRepo := NewUserRepository(db)

		user1 := createTestUserWithEmail(t, "duplicate_test@example.com")
		err := userRepo.Create(user1)
		require.NoError(t, err)

		user2, err := entity.NewUser("testuser2", "duplicate_test@example.com", "password456", "user")
		require.NoError(t, err)

		err = userRepo.Create(user2)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "email already exists")
	})

	t.Run("should return error for invalid user data", func(t *testing.T) {
		db := setupTestDB(t)
		userRepo := NewUserRepository(db)

		user := &entity.User{
			Username: "",
			Email:    "",
			Password: "",
			Role:     "",
		}

		err := userRepo.Create(user)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "email cannot be empty")
	})
}

func TestUser_FindByEmail(t *testing.T) {
	t.Run("should find user by email successfully", func(t *testing.T) {
		db := setupTestDB(t)
		userRepo := NewUserRepository(db)

		originalUser := createTestUserWithEmail(t, "find_test@example.com")
		err := userRepo.Create(originalUser)
		require.NoError(t, err)

		// Verify user was saved
		var count int64
		db.Model(&entity.User{}).Where("email = ?", "find_test@example.com").Count(&count)
		require.Equal(t, int64(1), count, "User should be saved before finding")

		foundUser, err := userRepo.FindByEmail("find_test@example.com")

		assert.NoError(t, err)
		assert.NotNil(t, foundUser)
		assert.Equal(t, originalUser.ID, foundUser.ID)
		assert.Equal(t, originalUser.Username, foundUser.Username)
		assert.Equal(t, originalUser.Email, foundUser.Email)
		assert.Equal(t, originalUser.Role, foundUser.Role)
		assert.NotEmpty(t, foundUser.Password)
	})

	t.Run("should return error when user not found", func(t *testing.T) {
		db := setupTestDB(t)
		userRepo := NewUserRepository(db)

		foundUser, err := userRepo.FindByEmail("nonexistent@example.com")

		assert.Error(t, err)
		assert.Nil(t, foundUser)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("should return error for empty email", func(t *testing.T) {
		db := setupTestDB(t)
		userRepo := NewUserRepository(db)

		foundUser, err := userRepo.FindByEmail("")

		assert.Error(t, err)
		assert.Nil(t, foundUser)
	})

	t.Run("should handle case insensitive email search", func(t *testing.T) {
		db := setupTestDB(t)
		userRepo := NewUserRepository(db)

		user := createTestUserWithEmail(t, "case_test@example.com")
		err := userRepo.Create(user)
		require.NoError(t, err)

		// Verify user was saved
		var count int64
		db.Model(&entity.User{}).Where("email = ?", "case_test@example.com").Count(&count)
		require.Equal(t, int64(1), count, "User should be saved before case insensitive search")

		foundUser, err := userRepo.FindByEmail("CASE_TEST@EXAMPLE.COM")

		assert.NoError(t, err)
		assert.NotNil(t, foundUser)
		assert.Equal(t, "case_test@example.com", foundUser.Email)
	})
}

func TestUser_Integration(t *testing.T) {
	t.Run("should create multiple users and find them", func(t *testing.T) {
		db := setupTestDB(t)
		userRepo := NewUserRepository(db)

		user1, err := entity.NewUser("user1", "integration_user1@example.com", "password1", "admin")
		require.NoError(t, err)

		user2, err := entity.NewUser("user2", "integration_user2@example.com", "password2", "user")
		require.NoError(t, err)

		err = userRepo.Create(user1)
		require.NoError(t, err)

		err = userRepo.Create(user2)
		require.NoError(t, err)

		// Verify both users were saved
		var count int64
		db.Model(&entity.User{}).Count(&count)
		require.Equal(t, int64(2), count, "Both users should be saved")

		foundUser1, err := userRepo.FindByEmail("integration_user1@example.com")
		assert.NoError(t, err)
		assert.Equal(t, "user1", foundUser1.Username)
		assert.Equal(t, "admin", foundUser1.Role)

		foundUser2, err := userRepo.FindByEmail("integration_user2@example.com")
		assert.NoError(t, err)
		assert.Equal(t, "user2", foundUser2.Username)
		assert.Equal(t, "user", foundUser2.Role)
	})
}

func BenchmarkUser_Create(b *testing.B) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&entity.User{})
	userRepo := NewUserRepository(db)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user, _ := entity.NewUser("benchuser", fmt.Sprintf("bench%d@example.com", i), "password", "user")
		userRepo.Create(user)
	}
}

func BenchmarkUser_FindByEmail(b *testing.B) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&entity.User{})
	userRepo := NewUserRepository(db)

	user, _ := entity.NewUser("testuser", "bench_find@example.com", "password123", "user")
	userRepo.Create(user)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		userRepo.FindByEmail("bench_find@example.com")
	}
}
