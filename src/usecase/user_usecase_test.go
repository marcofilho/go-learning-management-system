package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/infrastructure/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockTokenProvider for testing
type MockTokenProvider struct {
	mock.Mock
}

func (m *MockTokenProvider) GenerateToken(user *entity.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockTokenProvider) ValidateToken(tokenString string) (*auth.Claims, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.Claims), args.Error(1)
}

func TestUserUseCase_Register(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		password  string
		firstName string
		lastName  string
		role      entity.UserRole
		setupMock func(*MockUserRepository)
		wantErr   bool
		checkErr  func(error) bool
	}{
		{
			name:      "successful registration",
			email:     "test@example.com",
			password:  "password123",
			firstName: "John",
			lastName:  "Doe",
			role:      entity.UserRoleStudent,
			setupMock: func(repo *MockUserRepository) {
				repo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, entity.ErrNotFound)
				repo.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "duplicate email",
			email:     "existing@example.com",
			password:  "password123",
			firstName: "John",
			lastName:  "Doe",
			role:      entity.UserRoleStudent,
			setupMock: func(repo *MockUserRepository) {
				existingUser, _ := entity.NewUser("existing@example.com", "password123", "Existing", "User", entity.UserRoleStudent)
				repo.On("GetByEmail", mock.Anything, "existing@example.com").Return(existingUser, nil)
			},
			wantErr: true,
			checkErr: func(err error) bool {
				return errors.Is(err, entity.ErrDuplicateEntry)
			},
		},
		{
			name:      "repository error",
			email:     "test@example.com",
			password:  "password123",
			firstName: "John",
			lastName:  "Doe",
			role:      entity.UserRoleStudent,
			setupMock: func(repo *MockUserRepository) {
				repo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, entity.ErrNotFound)
				repo.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).Return(errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			mockTokenProvider := new(MockTokenProvider)
			tt.setupMock(mockRepo)

			uc := NewUserUseCase(mockRepo, mockTokenProvider)
			user, err := uc.Register(context.Background(), tt.email, tt.password, tt.firstName, tt.lastName, tt.role)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.checkErr != nil {
					assert.True(t, tt.checkErr(err))
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, user)
				assert.Equal(t, tt.email, user.Email)
				assert.Equal(t, tt.firstName, user.FirstName)
				assert.Equal(t, tt.lastName, user.LastName)
				assert.Equal(t, tt.role, user.Role)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserUseCase_Login(t *testing.T) {
	validUser, _ := entity.NewUser("test@example.com", "password123", "John", "Doe", entity.UserRoleStudent)

	tests := []struct {
		name      string
		email     string
		password  string
		setupMock func(*MockUserRepository, *MockTokenProvider)
		wantErr   bool
		checkErr  func(error) bool
	}{
		{
			name:     "successful login",
			email:    "test@example.com",
			password: "password123",
			setupMock: func(repo *MockUserRepository, tokenProvider *MockTokenProvider) {
				repo.On("GetByEmail", mock.Anything, "test@example.com").Return(validUser, nil)
				tokenProvider.On("GenerateToken", validUser).Return("valid-token", nil)
			},
			wantErr: false,
		},
		{
			name:     "user not found",
			email:    "notfound@example.com",
			password: "password123",
			setupMock: func(repo *MockUserRepository, tokenProvider *MockTokenProvider) {
				repo.On("GetByEmail", mock.Anything, "notfound@example.com").Return(nil, entity.ErrNotFound)
			},
			wantErr: true,
			checkErr: func(err error) bool {
				return errors.Is(err, entity.ErrInvalidCredentials)
			},
		},
		{
			name:     "wrong password",
			email:    "test@example.com",
			password: "wrongpassword",
			setupMock: func(repo *MockUserRepository, tokenProvider *MockTokenProvider) {
				repo.On("GetByEmail", mock.Anything, "test@example.com").Return(validUser, nil)
			},
			wantErr: true,
			checkErr: func(err error) bool {
				return errors.Is(err, entity.ErrInvalidCredentials)
			},
		},
		{
			name:     "inactive user",
			email:    "inactive@example.com",
			password: "password123",
			setupMock: func(repo *MockUserRepository, tokenProvider *MockTokenProvider) {
				inactiveUser, _ := entity.NewUser("inactive@example.com", "password123", "Inactive", "User", entity.UserRoleStudent)
				inactiveUser.IsActive = false
				repo.On("GetByEmail", mock.Anything, "inactive@example.com").Return(inactiveUser, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			mockTokenProvider := new(MockTokenProvider)
			tt.setupMock(mockRepo, mockTokenProvider)

			uc := NewUserUseCase(mockRepo, mockTokenProvider)
			token, user, err := uc.Login(context.Background(), tt.email, tt.password)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.checkErr != nil {
					assert.True(t, tt.checkErr(err))
				}
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, token)
				require.NotNil(t, user)
				assert.Equal(t, tt.email, user.Email)
			}

			mockRepo.AssertExpectations(t)
			mockTokenProvider.AssertExpectations(t)
		})
	}
}

func TestUserUseCase_GetUserByID(t *testing.T) {
	validUser, _ := entity.NewUser("test@example.com", "password123", "John", "Doe", entity.UserRoleStudent)

	mockRepo := new(MockUserRepository)
	mockTokenProvider := new(MockTokenProvider)
	mockRepo.On("GetByID", mock.Anything, validUser.ID).Return(validUser, nil)

	uc := NewUserUseCase(mockRepo, mockTokenProvider)
	user, err := uc.GetUserByID(context.Background(), validUser.ID)

	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, validUser.ID, user.ID)
	mockRepo.AssertExpectations(t)
}

func TestUserUseCase_ListUsers(t *testing.T) {
	user1, _ := entity.NewUser("user1@example.com", "password123", "User", "One", entity.UserRoleStudent)
	user2, _ := entity.NewUser("user2@example.com", "password123", "User", "Two", entity.UserRoleInstructor)
	users := []*entity.User{user1, user2}

	mockRepo := new(MockUserRepository)
	mockTokenProvider := new(MockTokenProvider)
	mockRepo.On("List", mock.Anything, 10, 0).Return(users, nil)

	uc := NewUserUseCase(mockRepo, mockTokenProvider)
	result, err := uc.ListUsers(context.Background(), 10, 0)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}

func TestUserUseCase_UpdateUser(t *testing.T) {
	user, _ := entity.NewUser("test@example.com", "password123", "John", "Doe", entity.UserRoleStudent)

	mockRepo := new(MockUserRepository)
	mockTokenProvider := new(MockTokenProvider)
	mockRepo.On("Update", mock.Anything, user).Return(nil)

	uc := NewUserUseCase(mockRepo, mockTokenProvider)
	err := uc.UpdateUser(context.Background(), user)

	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserUseCase_DeleteUser(t *testing.T) {
	userID := "test-user-id"

	mockRepo := new(MockUserRepository)
	mockTokenProvider := new(MockTokenProvider)
	mockRepo.On("Delete", mock.Anything, userID).Return(nil)

	uc := NewUserUseCase(mockRepo, mockTokenProvider)
	err := uc.DeleteUser(context.Background(), userID)

	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
