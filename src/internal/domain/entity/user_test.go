package entity

import (
"testing"
"github.com/stretchr/testify/assert"
"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	user, err := NewUser("test@example.com", "password123", "John", "Doe", UserRoleStudent)
	require.NoError(t, err)
	require.NotNil(t, user)
	assert.NotEmpty(t, user.ID)
	assert.Equal(t, "test@example.com", user.Email)
	assert.True(t, user.IsActive)
}

func TestUser_ValidatePassword(t *testing.T) {
	user, _ := NewUser("test@example.com", "correctpass", "John", "Doe", UserRoleStudent)
	assert.NoError(t, user.ValidatePassword("correctpass"))
	assert.Error(t, user.ValidatePassword("wrongpass"))
}

func TestUser_Roles(t *testing.T) {
	admin, err := NewUser("admin@example.com", "password123", "Admin", "User", UserRoleAdmin)
	require.NoError(t, err)
	assert.True(t, admin.IsAdmin())
	
	instructor, err := NewUser("inst@example.com", "password123", "Inst", "User", UserRoleInstructor)
	require.NoError(t, err)
	assert.True(t, instructor.IsInstructor())
	
	student, err := NewUser("student@example.com", "password123", "Student", "User", UserRoleStudent)
	require.NoError(t, err)
	assert.False(t, student.IsAdmin())
	assert.False(t, student.IsInstructor())
}
