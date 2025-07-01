package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var username = "testuser"
var email = "testuser@example.com"
var password = "password123"
var role = "user"

func TestNewUser(t *testing.T) {
	user, err := NewUser(username, email, password, role)
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, username, user.Username)
	assert.Equal(t, email, user.Email)
	assert.NotEmpty(t, user.Password) // Password should be hashed
	assert.Equal(t, role, user.Role)
	assert.NotEmpty(t, user.ID) // ID should be generated

}

func TestCheckPassword(t *testing.T) {
	user, err := NewUser(username, email, password, role)
	assert.Nil(t, err)

	// Check with correct password
	assert.True(t, user.CheckPassword(password))

	// Check with incorrect password
	assert.False(t, user.CheckPassword("wrongpassword"))
}
