package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/adapter/http/dto"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/adapter/http/middleware"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/infrastructure/auth"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

	req := httptest.NewRequest(http.MethodGet, "/users/"+user.ID.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"id": user.ID.String()})
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
	adminID := uuid.New()
	claims := &auth.Claims{
		UserID: adminID,
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
