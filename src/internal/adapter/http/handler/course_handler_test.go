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

	instructorID := uuid.New()
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

func TestCourseHandler_CreateCourse_Unauthorized(t *testing.T) {
	mockCourseRepo := new(MockCourseRepository)
	mockUserRepo := new(MockUserRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	courseUC := usecase.NewCourseUseCase(mockCourseRepo, mockUserRepo, mockAuditRepo)
	handler := NewCourseHandler(courseUC)

	reqBody := dto.CreateCourseRequest{Title: "Course", Description: "Desc"}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/courses", bytes.NewBuffer(body))

	rr := httptest.NewRecorder()
	handler.CreateCourse(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestCourseHandler_CreateCourse_InvalidJSON(t *testing.T) {
	mockCourseRepo := new(MockCourseRepository)
	mockUserRepo := new(MockUserRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	courseUC := usecase.NewCourseUseCase(mockCourseRepo, mockUserRepo, mockAuditRepo)
	handler := NewCourseHandler(courseUC)

	userID := uuid.New()
	req := httptest.NewRequest(http.MethodPost, "/courses", bytes.NewBuffer([]byte("{invalid")))
	claims := &auth.Claims{UserID: userID, Email: "inst@example.com", Role: entity.UserRoleInstructor}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, claims)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.CreateCourse(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCourseHandler_CreateCourse_AdminInvalidInstructorID(t *testing.T) {
	mockCourseRepo := new(MockCourseRepository)
	mockUserRepo := new(MockUserRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	courseUC := usecase.NewCourseUseCase(mockCourseRepo, mockUserRepo, mockAuditRepo)
	handler := NewCourseHandler(courseUC)

	adminID := uuid.New()
	reqBody := dto.CreateCourseRequest{Title: "Course", Description: "Desc", InstructorID: "not-a-uuid"}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/courses", bytes.NewBuffer(body))
	claims := &auth.Claims{UserID: adminID, Email: "admin@example.com", Role: entity.UserRoleAdmin}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, claims)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.CreateCourse(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCourseHandler_CreateCourse_InvalidDifficulty_MapsToBadRequest(t *testing.T) {
	mockCourseRepo := new(MockCourseRepository)
	mockUserRepo := new(MockUserRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	courseUC := usecase.NewCourseUseCase(mockCourseRepo, mockUserRepo, mockAuditRepo)
	handler := NewCourseHandler(courseUC)

	instructorID := uuid.New()
	instructor, _ := entity.NewUser("inst@example.com", "password123", "John", "Instructor", entity.UserRoleInstructor)
	instructor.ID = instructorID

	// user lookup ok
	mockUserRepo.On("GetByID", mock.Anything, instructorID).Return(instructor, nil)

	// entity.NewCourse will fail validation; due to current UC flow, Create is still invoked with nil
	mockCourseRepo.On("Create", mock.Anything, mock.Anything).Return(entity.ErrInvalidInput)

	reqBody := dto.CreateCourseRequest{Title: "Course", Description: "Desc", DifficultyLevel: "invalid-level"}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/courses", bytes.NewBuffer(body))
	claims := &auth.Claims{UserID: instructorID, Email: "inst@example.com", Role: entity.UserRoleInstructor}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, claims)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.CreateCourse(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockUserRepo.AssertExpectations(t)
}

func TestCourseHandler_CreateCourse_AdminForOtherInstructor_Success(t *testing.T) {
	mockCourseRepo := new(MockCourseRepository)
	mockUserRepo := new(MockUserRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	courseUC := usecase.NewCourseUseCase(mockCourseRepo, mockUserRepo, mockAuditRepo)
	handler := NewCourseHandler(courseUC)

	adminID := uuid.New()
	instructorID := uuid.New()
	instructor, _ := entity.NewUser("inst@example.com", "password123", "John", "Instructor", entity.UserRoleInstructor)
	instructor.ID = instructorID

	mockUserRepo.On("GetByID", mock.Anything, instructorID).Return(instructor, nil)
	mockCourseRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Course")).Return(nil)
	mockAuditRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.AuditLog")).Return(nil)

	reqBody := dto.CreateCourseRequest{
		Title:        "Test Course",
		Description:  "Test Description",
		InstructorID: instructorID.String(),
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/courses", bytes.NewBuffer(body))

	claims := &auth.Claims{UserID: adminID, Email: "admin@example.com", Role: entity.UserRoleAdmin}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, claims)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.CreateCourse(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	mockCourseRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestCourseHandler_CreateCourse_NonAdminSettingInstructor_Forbidden(t *testing.T) {
	mockCourseRepo := new(MockCourseRepository)
	mockUserRepo := new(MockUserRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	courseUC := usecase.NewCourseUseCase(mockCourseRepo, mockUserRepo, mockAuditRepo)
	handler := NewCourseHandler(courseUC)

	userID := uuid.New()
	otherInstructorID := uuid.New()

	reqBody := dto.CreateCourseRequest{
		Title:        "Course",
		InstructorID: otherInstructorID.String(),
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/courses", bytes.NewBuffer(body))

	claims := &auth.Claims{UserID: userID, Email: "inst@example.com", Role: entity.UserRoleInstructor}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, claims)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.CreateCourse(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
}

func TestCourseHandler_ListCourses_WithFilters(t *testing.T) {
	mockCourseRepo := new(MockCourseRepository)
	mockUserRepo := new(MockUserRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	courseUC := usecase.NewCourseUseCase(mockCourseRepo, mockUserRepo, mockAuditRepo)
	handler := NewCourseHandler(courseUC)

	instructorID := uuid.New()
	courses := []*entity.Course{
		{ID: uuid.New(), Title: "Course 1", InstructorID: instructorID, DifficultyLevel: entity.DifficultyLevelIntermediate},
	}

	mockCourseRepo.On("List", mock.Anything, mock.AnythingOfType("*repository.CourseFilter")).Return(courses, nil)

	url := "/courses?instructor_id=" + instructorID.String() + "&difficulty_level=intermediate&active_only=true&limit=5&offset=10"
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rr := httptest.NewRecorder()

	handler.ListCourses(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockCourseRepo.AssertExpectations(t)
}

func TestCourseHandler_GetCourse_Success(t *testing.T) {
	mockCourseRepo := new(MockCourseRepository)
	mockUserRepo := new(MockUserRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	courseUC := usecase.NewCourseUseCase(mockCourseRepo, mockUserRepo, mockAuditRepo)
	handler := NewCourseHandler(courseUC)

	courseID := uuid.New()
	course := &entity.Course{
		ID:              courseID,
		Title:           "Test Course",
		Description:     "Test Description",
		InstructorID:    uuid.New(),
		DifficultyLevel: entity.DifficultyLevelBeginner,
	}

	mockCourseRepo.On("GetByID", mock.Anything, courseID).Return(course, nil)

	req := httptest.NewRequest(http.MethodGet, "/courses/"+courseID.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"id": courseID.String()})
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

	courseID := uuid.New()
	mockCourseRepo.On("GetByID", mock.Anything, courseID).Return(nil, entity.ErrNotFound)

	req := httptest.NewRequest(http.MethodGet, "/courses/"+courseID.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"id": courseID.String()})
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
		{ID: uuid.New(), Title: "Course 1"},
		{ID: uuid.New(), Title: "Course 2"},
	}

	mockCourseRepo.On("List", mock.Anything, mock.AnythingOfType("*repository.CourseFilter")).Return(courses, nil)

	req := httptest.NewRequest(http.MethodGet, "/courses", nil)
	rr := httptest.NewRecorder()

	handler.ListCourses(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockCourseRepo.AssertExpectations(t)
}

func TestCourseHandler_ListCourses_InvalidInstructorID(t *testing.T) {
	mockCourseRepo := new(MockCourseRepository)
	mockUserRepo := new(MockUserRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	courseUC := usecase.NewCourseUseCase(mockCourseRepo, mockUserRepo, mockAuditRepo)
	handler := NewCourseHandler(courseUC)

	req := httptest.NewRequest(http.MethodGet, "/courses?instructor_id=not-a-uuid", nil)
	rr := httptest.NewRecorder()

	handler.ListCourses(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCourseHandler_DeleteCourse_Success(t *testing.T) {
	mockCourseRepo := new(MockCourseRepository)
	mockUserRepo := new(MockUserRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	courseUC := usecase.NewCourseUseCase(mockCourseRepo, mockUserRepo, mockAuditRepo)
	handler := NewCourseHandler(courseUC)

	courseID := uuid.New()
	instructorID := uuid.New()
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

	req := httptest.NewRequest(http.MethodDelete, "/courses/"+courseID.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"id": courseID.String()})

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

func TestCourseHandler_DeleteCourse_Unauthorized(t *testing.T) {
	mockCourseRepo := new(MockCourseRepository)
	mockUserRepo := new(MockUserRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	courseUC := usecase.NewCourseUseCase(mockCourseRepo, mockUserRepo, mockAuditRepo)
	handler := NewCourseHandler(courseUC)

	courseID := uuid.New()

	req := httptest.NewRequest(http.MethodDelete, "/courses/"+courseID.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"id": courseID.String()})

	rr := httptest.NewRecorder()
	handler.DeleteCourse(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestCourseHandler_DeleteCourse_InvalidID(t *testing.T) {
	mockCourseRepo := new(MockCourseRepository)
	mockUserRepo := new(MockUserRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	courseUC := usecase.NewCourseUseCase(mockCourseRepo, mockUserRepo, mockAuditRepo)
	handler := NewCourseHandler(courseUC)

	req := httptest.NewRequest(http.MethodDelete, "/courses/invalid-uuid", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "invalid-uuid"})

	rr := httptest.NewRecorder()
	handler.DeleteCourse(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCourseHandler_UpdateCourse_Success(t *testing.T) {
	mockCourseRepo := new(MockCourseRepository)
	mockUserRepo := new(MockUserRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	courseUC := usecase.NewCourseUseCase(mockCourseRepo, mockUserRepo, mockAuditRepo)
	handler := NewCourseHandler(courseUC)

	courseID := uuid.New()
	instructorID := uuid.New()
	instructor, _ := entity.NewUser("inst@example.com", "password123", "John", "Instructor", entity.UserRoleInstructor)
	instructor.ID = instructorID
	course := &entity.Course{
		ID:           courseID,
		Title:        "Old Title",
		InstructorID: instructorID,
	}

	mockUserRepo.On("GetByID", mock.Anything, instructorID).Return(instructor, nil)
	mockCourseRepo.On("GetByID", mock.Anything, courseID).Return(course, nil)
	mockCourseRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Course")).Return(nil)
	mockAuditRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.AuditLog")).Return(nil)

	reqBody := dto.UpdateCourseRequest{
		Title:           "New Title",
		Description:     "New Description",
		DifficultyLevel: string(entity.DifficultyLevelIntermediate),
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/courses/"+courseID.String(), bytes.NewBuffer(body))
	req = mux.SetURLVars(req, map[string]string{"id": courseID.String()})

	claims := &auth.Claims{
		UserID: instructorID,
		Email:  "inst@example.com",
		Role:   entity.UserRoleInstructor,
	}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, claims)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.UpdateCourse(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockCourseRepo.AssertExpectations(t)
}

func TestCourseHandler_UpdateCourse_Forbidden_NotInstructor(t *testing.T) {
	mockCourseRepo := new(MockCourseRepository)
	mockUserRepo := new(MockUserRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	courseUC := usecase.NewCourseUseCase(mockCourseRepo, mockUserRepo, mockAuditRepo)
	handler := NewCourseHandler(courseUC)

	courseID := uuid.New()
	instructorID := uuid.New()
	otherUserID := uuid.New()
	instructor, _ := entity.NewUser("inst@example.com", "password123", "John", "Instructor", entity.UserRoleInstructor)
	instructor.ID = instructorID
	course := &entity.Course{
		ID:           courseID,
		Title:        "Old Title",
		InstructorID: instructorID,
	}

	mockUserRepo.On("GetByID", mock.Anything, instructorID).Return(instructor, nil)
	mockCourseRepo.On("GetByID", mock.Anything, courseID).Return(course, nil)

	reqBody := dto.UpdateCourseRequest{Title: "New Title"}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/courses/"+courseID.String(), bytes.NewBuffer(body))
	req = mux.SetURLVars(req, map[string]string{"id": courseID.String()})

	claims := &auth.Claims{UserID: otherUserID, Email: "other@example.com", Role: entity.UserRoleInstructor}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, claims)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.UpdateCourse(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
}

func TestCourseHandler_UpdateCourse_InvalidPayload(t *testing.T) {
	mockCourseRepo := new(MockCourseRepository)
	mockUserRepo := new(MockUserRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	courseUC := usecase.NewCourseUseCase(mockCourseRepo, mockUserRepo, mockAuditRepo)
	handler := NewCourseHandler(courseUC)

	courseID := uuid.New()
	instructorID := uuid.New()
	instructor, _ := entity.NewUser("inst@example.com", "password123", "John", "Instructor", entity.UserRoleInstructor)
	instructor.ID = instructorID
	course := &entity.Course{
		ID:           courseID,
		Title:        "Old Title",
		InstructorID: instructorID,
	}

	mockUserRepo.On("GetByID", mock.Anything, instructorID).Return(instructor, nil)
	mockCourseRepo.On("GetByID", mock.Anything, courseID).Return(course, nil)

	req := httptest.NewRequest(http.MethodPut, "/courses/"+courseID.String(), bytes.NewBuffer([]byte("{invalid json")))
	req = mux.SetURLVars(req, map[string]string{"id": courseID.String()})

	claims := &auth.Claims{UserID: instructorID, Email: "inst@example.com", Role: entity.UserRoleInstructor}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, claims)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.UpdateCourse(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
