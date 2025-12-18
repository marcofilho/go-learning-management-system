package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/adapter/http/dto"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/adapter/http/middleware"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/infrastructure/auth"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetByRole(ctx context.Context, role entity.UserRole) ([]*entity.User, error) {
	args := m.Called(ctx, role)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (m *MockUserRepository) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

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

func TestUserHandler_Register_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockToken := new(MockTokenProvider)
	userUC := usecase.NewUserUseCase(mockRepo, mockToken)
	handler := NewUserHandler(userUC)

	mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, entity.ErrNotFound)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil)

	reqBody := dto.RegisterRequest{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	handler.Register(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	mockRepo.AssertExpectations(t)
}

func TestUserHandler_Register_DuplicateEmail(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockToken := new(MockTokenProvider)
	userUC := usecase.NewUserUseCase(mockRepo, mockToken)
	handler := NewUserHandler(userUC)

	existingUser, _ := entity.NewUser("test@example.com", "password123", "Existing", "User", entity.UserRoleStudent)
	mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(existingUser, nil)

	reqBody := dto.RegisterRequest{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	handler.Register(rr, req)

	assert.Equal(t, http.StatusConflict, rr.Code)
	mockRepo.AssertExpectations(t)
}

func TestUserHandler_Login_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockToken := new(MockTokenProvider)
	userUC := usecase.NewUserUseCase(mockRepo, mockToken)
	handler := NewUserHandler(userUC)

	user, _ := entity.NewUser("test@example.com", "password123", "John", "Doe", entity.UserRoleStudent)
	mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(user, nil)
	mockToken.On("GenerateToken", user).Return("token123", nil)

	reqBody := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	handler.Login(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockRepo.AssertExpectations(t)
	mockToken.AssertExpectations(t)
}

func TestUserHandler_GetUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockToken := new(MockTokenProvider)
	userUC := usecase.NewUserUseCase(mockRepo, mockToken)
	handler := NewUserHandler(userUC)

	user, _ := entity.NewUser("test@example.com", "password123", "John", "Doe", entity.UserRoleStudent)
	mockRepo.On("GetByID", mock.Anything, user.ID).Return(user, nil)

	req := httptest.NewRequest(http.MethodGet, "/users/"+user.ID, nil)
	req = mux.SetURLVars(req, map[string]string{"id": user.ID})
	rr := httptest.NewRecorder()

	handler.GetUserByID(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockRepo.AssertExpectations(t)
}

func TestUserHandler_ListUsers_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockToken := new(MockTokenProvider)
	userUC := usecase.NewUserUseCase(mockRepo, mockToken)
	handler := NewUserHandler(userUC)

	user1, _ := entity.NewUser("test1@example.com", "password123", "John", "Doe", entity.UserRoleStudent)
	user2, _ := entity.NewUser("test2@example.com", "password123", "Jane", "Smith", entity.UserRoleInstructor)
	users := []*entity.User{user1, user2}
	
	mockRepo.On("List", mock.Anything, 10, 0).Return(users, nil)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	claims := &auth.Claims{
		UserID: "admin-id",
		Email:  "admin@example.com",
		Role:   entity.UserRoleAdmin,
	}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, claims)
	req = req.WithContext(ctx)
	
	rr := httptest.NewRecorder()

	handler.ListUsers(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockRepo.AssertExpectations(t)
}
