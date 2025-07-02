package database

import (
	"errors"
	"strings"
	"time"

	"github/GuilhermeHermes/GO_API/internal/entity"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (u *UserRepository) Create(user *entity.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}

	if strings.TrimSpace(user.Email) == "" {
		return errors.New("email cannot be empty")
	}

	if strings.TrimSpace(user.Username) == "" {
		return errors.New("username cannot be empty")
	}

	// Normalize email for consistent storage and comparison
	normalizedEmail := strings.ToLower(strings.TrimSpace(user.Email))
	user.Email = normalizedEmail

	var existingUser entity.User
	err := u.DB.Where("email = ?", normalizedEmail).First(&existingUser).Error
	if err == nil {
		return errors.New("email already exists")
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}

	// Set timestamps manually since entity uses string fields
	now := time.Now().Format(time.RFC3339)
	if user.CreatedAt == "" {
		user.CreatedAt = now
	}
	user.UpdatedAt = now

	return u.DB.Create(user).Error
}

func (u *UserRepository) FindByEmail(email string) (*entity.User, error) {
	if strings.TrimSpace(email) == "" {
		return nil, errors.New("email cannot be empty")
	}

	normalizedEmail := strings.ToLower(strings.TrimSpace(email))

	var user entity.User
	if err := u.DB.Where("email = ?", normalizedEmail).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserRepository) FindByID(id string) (*entity.User, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("id cannot be empty")
	}

	var user entity.User
	if err := u.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserRepository) Update(user *entity.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}

	// Set updated timestamp manually
	user.UpdatedAt = time.Now().Format(time.RFC3339)
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))

	return u.DB.Save(user).Error
}

func (u *UserRepository) Delete(id string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("id cannot be empty")
	}

	return u.DB.Delete(&entity.User{}, "id = ?", id).Error
}

func (u *UserRepository) Exists(email string) (bool, error) {
	if strings.TrimSpace(email) == "" {
		return false, errors.New("email cannot be empty")
	}

	// Normalize email for consistent comparison
	normalizedEmail := strings.ToLower(strings.TrimSpace(email))

	var count int64
	err := u.DB.Model(&entity.User{}).Where("email = ?", normalizedEmail).Count(&count).Error
	return count > 0, err
}
