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
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/repository"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/infrastructure/auth"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEnrollmentHandler_EnrollInCourse_SelfEnrollment(t *testing.T) {
	mockEnrollRepo := new(MockEnrollmentRepository)
	mockUserRepo := new(MockUserRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	enrollmentUC := usecase.NewEnrollmentUseCase(mockEnrollRepo, mockCourseRepo, mockUserRepo, mockAuditRepo)
	handler := NewEnrollmentHandler(enrollmentUC)

	studentID := uuid.New()
	courseID := uuid.New()
	student, _ := entity.NewUser("student@example.com", "password123", "John", "Student", entity.UserRoleStudent)
	student.ID = studentID
	course := &entity.Course{ID: courseID, Title: "Test Course"}

	mockUserRepo.On("GetByID", mock.Anything, studentID).Return(student, nil)
	mockCourseRepo.On("GetByID", mock.Anything, courseID).Return(course, nil)
	mockEnrollRepo.On("GetByStudentAndCourse", mock.Anything, studentID, courseID).Return(nil, entity.ErrNotFound)
	mockEnrollRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Enrollment")).Return(nil)
	mockAuditRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.AuditLog")).Return(nil)

	reqBody := dto.EnrollCourseRequest{}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/courses/"+courseID.String()+"/enroll", bytes.NewBuffer(body))
	req = mux.SetURLVars(req, map[string]string{"id": courseID.String()})

	claims := &auth.Claims{
		UserID: studentID,
		Email:  "student@example.com",
		Role:   entity.UserRoleStudent,
	}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, claims)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.EnrollInCourse(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	mockEnrollRepo.AssertExpectations(t)
}

func TestEnrollmentHandler_EnrollInCourse_InvalidCourseID(t *testing.T) {
	mockEnrollRepo := new(MockEnrollmentRepository)
	mockUserRepo := new(MockUserRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	enrollmentUC := usecase.NewEnrollmentUseCase(mockEnrollRepo, mockCourseRepo, mockUserRepo, mockAuditRepo)
	handler := NewEnrollmentHandler(enrollmentUC)

	studentID := uuid.New()

	req := httptest.NewRequest(http.MethodPost, "/courses/invalid-uuid/enroll", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "invalid-uuid"})

	claims := &auth.Claims{UserID: studentID, Email: "student@example.com", Role: entity.UserRoleStudent}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, claims)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.EnrollInCourse(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestEnrollmentHandler_EnrollInCourse_Unauthorized(t *testing.T) {
	mockEnrollRepo := new(MockEnrollmentRepository)
	mockUserRepo := new(MockUserRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	enrollmentUC := usecase.NewEnrollmentUseCase(mockEnrollRepo, mockCourseRepo, mockUserRepo, mockAuditRepo)
	handler := NewEnrollmentHandler(enrollmentUC)

	courseID := uuid.New()
	req := httptest.NewRequest(http.MethodPost, "/courses/"+courseID.String()+"/enroll", nil)
	req = mux.SetURLVars(req, map[string]string{"id": courseID.String()})

	rr := httptest.NewRecorder()
	handler.EnrollInCourse(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestEnrollmentHandler_EnrollInCourse_InvalidStudentIDInBody(t *testing.T) {
	mockEnrollRepo := new(MockEnrollmentRepository)
	mockUserRepo := new(MockUserRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	enrollmentUC := usecase.NewEnrollmentUseCase(mockEnrollRepo, mockCourseRepo, mockUserRepo, mockAuditRepo)
	handler := NewEnrollmentHandler(enrollmentUC)

	studentID := uuid.New()
	courseID := uuid.New()

	reqBody := dto.EnrollCourseRequest{StudentID: "not-a-uuid"}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/courses/"+courseID.String()+"/enroll", bytes.NewBuffer(body))
	req = mux.SetURLVars(req, map[string]string{"id": courseID.String()})

	claims := &auth.Claims{UserID: studentID, Email: "student@example.com", Role: entity.UserRoleStudent}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, claims)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.EnrollInCourse(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestEnrollmentHandler_GetStudentCourses_Success(t *testing.T) {
	mockEnrollRepo := new(MockEnrollmentRepository)
	mockUserRepo := new(MockUserRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	enrollmentUC := usecase.NewEnrollmentUseCase(mockEnrollRepo, mockCourseRepo, mockUserRepo, mockAuditRepo)
	handler := NewEnrollmentHandler(enrollmentUC)

	studentID := uuid.New()
	enrollments := []*entity.Enrollment{
		{StudentID: studentID, CourseID: uuid.New(), Status: entity.EnrollmentStatusActive},
		{StudentID: studentID, CourseID: uuid.New(), Status: entity.EnrollmentStatusActive},
	}

	mockEnrollRepo.On("GetByStudent", mock.Anything, studentID, mock.AnythingOfType("*repository.EnrollmentFilter")).Return(enrollments, int64(len(enrollments)), nil)

	req := httptest.NewRequest(http.MethodGet, "/students/"+studentID.String()+"/courses", nil)
	req = mux.SetURLVars(req, map[string]string{"id": studentID.String()})
	rr := httptest.NewRecorder()

	handler.GetStudentCourses(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockEnrollRepo.AssertExpectations(t)
}

func TestEnrollmentHandler_GetStudentCourses_WithFilters(t *testing.T) {
	mockEnrollRepo := new(MockEnrollmentRepository)
	mockUserRepo := new(MockUserRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	enrollmentUC := usecase.NewEnrollmentUseCase(mockEnrollRepo, mockCourseRepo, mockUserRepo, mockAuditRepo)
	handler := NewEnrollmentHandler(enrollmentUC)

	studentID := uuid.New()
	enrollments := []*entity.Enrollment{
		{StudentID: studentID, CourseID: uuid.New(), Status: entity.EnrollmentStatusActive},
	}

	mockEnrollRepo.On("GetByStudent", mock.Anything, studentID, mock.AnythingOfType("*repository.EnrollmentFilter")).Return(enrollments, int64(len(enrollments)), nil)

	url := "/students/" + studentID.String() + "/courses?status=active&date_from=2025-01-01&date_to=2025-12-31&limit=5&offset=2"
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req = mux.SetURLVars(req, map[string]string{"id": studentID.String()})
	rr := httptest.NewRecorder()

	handler.GetStudentCourses(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockEnrollRepo.AssertExpectations(t)
}

func TestEnrollmentHandler_GetStudentCourses_InvalidStudentID(t *testing.T) {
	mockEnrollRepo := new(MockEnrollmentRepository)
	mockUserRepo := new(MockUserRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	enrollmentUC := usecase.NewEnrollmentUseCase(mockEnrollRepo, mockCourseRepo, mockUserRepo, mockAuditRepo)
	handler := NewEnrollmentHandler(enrollmentUC)

	req := httptest.NewRequest(http.MethodGet, "/students/invalid-uuid/courses", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "invalid-uuid"})
	rr := httptest.NewRecorder()

	handler.GetStudentCourses(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestEnrollmentHandler_GetCourseStudents_Success(t *testing.T) {
	mockEnrollRepo := new(MockEnrollmentRepository)
	mockUserRepo := new(MockUserRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	enrollmentUC := usecase.NewEnrollmentUseCase(mockEnrollRepo, mockCourseRepo, mockUserRepo, mockAuditRepo)
	handler := NewEnrollmentHandler(enrollmentUC)

	courseID := uuid.New()
	enrollments := []*entity.Enrollment{
		{StudentID: uuid.New(), CourseID: courseID, Status: entity.EnrollmentStatusActive},
		{StudentID: uuid.New(), CourseID: courseID, Status: entity.EnrollmentStatusActive},
	}

	mockEnrollRepo.On("GetByCourse", mock.Anything, courseID, mock.AnythingOfType("*repository.EnrollmentFilter")).Return(enrollments, int64(len(enrollments)), nil)

	req := httptest.NewRequest(http.MethodGet, "/courses/"+courseID.String()+"/students", nil)
	req = mux.SetURLVars(req, map[string]string{"id": courseID.String()})
	rr := httptest.NewRecorder()

	handler.GetCourseStudents(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockEnrollRepo.AssertExpectations(t)
}

func TestEnrollmentHandler_GetCourseStudents_InvalidCourseID(t *testing.T) {
	mockEnrollRepo := new(MockEnrollmentRepository)
	mockUserRepo := new(MockUserRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	enrollmentUC := usecase.NewEnrollmentUseCase(mockEnrollRepo, mockCourseRepo, mockUserRepo, mockAuditRepo)
	handler := NewEnrollmentHandler(enrollmentUC)

	req := httptest.NewRequest(http.MethodGet, "/courses/invalid-uuid/students", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "invalid-uuid"})
	rr := httptest.NewRecorder()

	handler.GetCourseStudents(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestEnrollmentHandler_GetStudentCourses_DefaultLimitOffsetOnInvalidValues(t *testing.T) {
	mockEnrollRepo := new(MockEnrollmentRepository)
	mockUserRepo := new(MockUserRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	enrollmentUC := usecase.NewEnrollmentUseCase(mockEnrollRepo, mockCourseRepo, mockUserRepo, mockAuditRepo)
	handler := NewEnrollmentHandler(enrollmentUC)

	studentID := uuid.New()
	matcher := mock.MatchedBy(func(f *repository.EnrollmentFilter) bool {
		return f.Limit == 10 && f.Offset == 0 && f.Status == nil && f.DateFrom == nil && f.DateTo == nil
	})
	mockEnrollRepo.On("GetByStudent", mock.Anything, studentID, matcher).Return([]*entity.Enrollment{}, int64(0), nil)

	req := httptest.NewRequest(http.MethodGet, "/students/"+studentID.String()+"/courses?limit=0&offset=-1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": studentID.String()})
	rr := httptest.NewRecorder()

	handler.GetStudentCourses(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockEnrollRepo.AssertExpectations(t)
}

func TestEnrollmentHandler_GetCourseStudents_StatusOnlyFilter(t *testing.T) {
	mockEnrollRepo := new(MockEnrollmentRepository)
	mockUserRepo := new(MockUserRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	enrollmentUC := usecase.NewEnrollmentUseCase(mockEnrollRepo, mockCourseRepo, mockUserRepo, mockAuditRepo)
	handler := NewEnrollmentHandler(enrollmentUC)

	courseID := uuid.New()
	statusMatcher := mock.MatchedBy(func(f *repository.EnrollmentFilter) bool {
		return f.Limit == 10 && f.Offset == 0 && f.Status != nil && *f.Status == entity.EnrollmentStatusActive && f.DateFrom == nil && f.DateTo == nil
	})
	mockEnrollRepo.On("GetByCourse", mock.Anything, courseID, statusMatcher).Return([]*entity.Enrollment{}, int64(0), nil)

	req := httptest.NewRequest(http.MethodGet, "/courses/"+courseID.String()+"/students?status=active", nil)
	req = mux.SetURLVars(req, map[string]string{"id": courseID.String()})
	rr := httptest.NewRecorder()

	handler.GetCourseStudents(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockEnrollRepo.AssertExpectations(t)
}
