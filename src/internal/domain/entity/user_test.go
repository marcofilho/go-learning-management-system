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
func TestUser_Validate(t *testing.T) {
	tests := []struct {
		name    string
		user    *User
		wantErr bool
	}{
		{
			name: "EmptyEmail",
			user: &User{
				Email:     "",
				FirstName: "John",
				LastName:  "Doe",
				Role:      UserRoleStudent,
			},
			wantErr: true,
		},
		{
			name: "InvalidEmail",
			user: &User{
				Email:     "invalid-email",
				FirstName: "John",
				LastName:  "Doe",
				Role:      UserRoleStudent,
			},
			wantErr: true,
		},
		{
			name: "EmptyFirstName",
			user: &User{
				Email:     "test@example.com",
				FirstName: "",
				LastName:  "Doe",
				Role:      UserRoleStudent,
			},
			wantErr: true,
		},
		{
			name: "EmptyLastName",
			user: &User{
				Email:     "test@example.com",
				FirstName: "John",
				LastName:  "",
				Role:      UserRoleStudent,
			},
			wantErr: true,
		},
		{
			name: "InvalidRole",
			user: &User{
				Email:     "test@example.com",
				FirstName: "John",
				LastName:  "Doe",
				Role:      "invalid-role",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNewUser_Error(t *testing.T) {
	_, err := NewUser("not-an-email", "password123", "John", "Doe", UserRoleStudent)
	assert.Error(t, err, "should fail with invalid email")

	// Note: Password length validation happens at usecase/handler level, not in entity
	// The entity layer accepts any password and hashes it with bcrypt
}

func TestUser_CanAuthenticate(t *testing.T) {
	user, _ := NewUser("test@example.com", "password123", "John", "Doe", UserRoleStudent)
	assert.NoError(t, user.CanAuthenticate())

	user.IsActive = false
	assert.Error(t, user.CanAuthenticate())
}

func TestUser_FullName(t *testing.T) {
	user := &User{FirstName: "John", LastName: "Doe"}
	assert.Equal(t, "John Doe", user.FullName())
}

func TestUser_TableName(t *testing.T) {
	var user User
	assert.Equal(t, "users", user.TableName())
}
