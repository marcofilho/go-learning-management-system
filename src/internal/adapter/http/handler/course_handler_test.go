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

func TestCourseHandler_CreateCourse_Success(t *testing.T) {
	mockCourseRepo := new(MockCourseRepository)
	mockUserRepo := new(MockUserRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	courseUC := usecase.NewCourseUseCase(mockCourseRepo, mockUserRepo, mockAuditRepo)
	handler := NewCourseHandler(courseUC)

	instructorID := uuid.New().String()
	instructor, _ := entity.NewUser("inst@example.com", "password123", "John", "Instructor", entity.UserRoleInstructor)
	instructor.ID = instructorID

	mockUserRepo.On("GetByID", mock.Anything, instructorID).Return(instructor, nil)
	mockCourseRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Course")).Return(nil)
	mockAuditRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.AuditLog")).Return(nil)

	reqBody := dto.CreateCourseRequest{
		Title:           "Test Course",
		Description:     "Test Description",
		DifficultyLevel: string(entity.DifficultyLevelBeginner),
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/courses", bytes.NewBuffer(body))

	claims := &auth.Claims{
		UserID: instructorID,
		Email:  "inst@example.com",
		Role:   entity.UserRoleInstructor,
	}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, claims)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.CreateCourse(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	mockCourseRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestCourseHandler_GetCourse_Success(t *testing.T) {
	mockCourseRepo := new(MockCourseRepository)
	mockUserRepo := new(MockUserRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	courseUC := usecase.NewCourseUseCase(mockCourseRepo, mockUserRepo, mockAuditRepo)
	handler := NewCourseHandler(courseUC)

	courseID := uuid.New().String()
	course := &entity.Course{
		ID:              courseID,
		Title:           "Test Course",
		Description:     "Test Description",
		InstructorID:    uuid.New().String(),
		DifficultyLevel: entity.DifficultyLevelBeginner,
	}

	mockCourseRepo.On("GetByID", mock.Anything, courseID).Return(course, nil)

	req := httptest.NewRequest(http.MethodGet, "/courses/"+courseID, nil)
	req = mux.SetURLVars(req, map[string]string{"id": courseID})
	rr := httptest.NewRecorder()

	handler.GetCourse(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockCourseRepo.AssertExpectations(t)
}

func TestCourseHandler_GetCourse_NotFound(t *testing.T) {
	mockCourseRepo := new(MockCourseRepository)
	mockUserRepo := new(MockUserRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	courseUC := usecase.NewCourseUseCase(mockCourseRepo, mockUserRepo, mockAuditRepo)
	handler := NewCourseHandler(courseUC)

	courseID := uuid.New().String()
	mockCourseRepo.On("GetByID", mock.Anything, courseID).Return(nil, entity.ErrNotFound)

	req := httptest.NewRequest(http.MethodGet, "/courses/"+courseID, nil)
	req = mux.SetURLVars(req, map[string]string{"id": courseID})
	rr := httptest.NewRecorder()

	handler.GetCourse(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	mockCourseRepo.AssertExpectations(t)
}

func TestCourseHandler_ListCourses_Success(t *testing.T) {
	mockCourseRepo := new(MockCourseRepository)
	mockUserRepo := new(MockUserRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	courseUC := usecase.NewCourseUseCase(mockCourseRepo, mockUserRepo, mockAuditRepo)
	handler := NewCourseHandler(courseUC)

	courses := []*entity.Course{
		{ID: uuid.New().String(), Title: "Course 1"},
		{ID: uuid.New().String(), Title: "Course 2"},
	}

	mockCourseRepo.On("List", mock.Anything, mock.AnythingOfType("*repository.CourseFilter")).Return(courses, nil)

	req := httptest.NewRequest(http.MethodGet, "/courses", nil)
	rr := httptest.NewRecorder()

	handler.ListCourses(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockCourseRepo.AssertExpectations(t)
}

func TestCourseHandler_DeleteCourse_Success(t *testing.T) {
	mockCourseRepo := new(MockCourseRepository)
	mockUserRepo := new(MockUserRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	courseUC := usecase.NewCourseUseCase(mockCourseRepo, mockUserRepo, mockAuditRepo)
	handler := NewCourseHandler(courseUC)

	courseID := uuid.New().String()
	instructorID := uuid.New().String()
	instructor, _ := entity.NewUser("inst@example.com", "password123", "John", "Instructor", entity.UserRoleInstructor)
	instructor.ID = instructorID
	course := &entity.Course{
		ID:           courseID,
		Title:        "Test Course",
		InstructorID: instructorID,
	}

	mockUserRepo.On("GetByID", mock.Anything, instructorID).Return(instructor, nil)
	mockCourseRepo.On("GetByID", mock.Anything, courseID).Return(course, nil)
	mockCourseRepo.On("Delete", mock.Anything, courseID).Return(nil)
	mockAuditRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.AuditLog")).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/courses/"+courseID, nil)
	req = mux.SetURLVars(req, map[string]string{"id": courseID})

	claims := &auth.Claims{
		UserID: instructorID,
		Email:  "inst@example.com",
		Role:   entity.UserRoleInstructor,
	}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, claims)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.DeleteCourse(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
	mockCourseRepo.AssertExpectations(t)
}
