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

func TestLessonHandler_CreateLesson_Success(t *testing.T) {
	mockLessonRepo := new(MockLessonRepository)
	mockModuleRepo := new(MockModuleRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	lessonUC := usecase.NewLessonUseCase(mockLessonRepo, mockModuleRepo, mockCourseRepo, mockAuditRepo)
	handler := NewLessonHandler(lessonUC)

	moduleID := uuid.New()
	instructorID := uuid.New()
	courseID := uuid.New()
	module := &entity.Module{ID: moduleID, CourseID: courseID}
	course := &entity.Course{ID: courseID, InstructorID: instructorID}

	mockModuleRepo.On("GetByID", mock.Anything, moduleID).Return(module, nil)
	mockCourseRepo.On("GetByID", mock.Anything, courseID).Return(course, nil)
	mockLessonRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.LessonVersion")).Return(nil)
	mockAuditRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.AuditLog")).Return(nil)

	reqBody := dto.CreateLessonRequest{
		Content:  "Lesson content",
		VideoURL: "https://example.com/video.mp4",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/modules/"+moduleID.String()+"/lessons", bytes.NewBuffer(body))
	req = mux.SetURLVars(req, map[string]string{"moduleId": moduleID.String()})

	claims := &auth.Claims{
		UserID: instructorID,
		Email:  "instructor@example.com",
		Role:   entity.UserRoleInstructor,
	}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, claims)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.CreateLesson(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	mockLessonRepo.AssertExpectations(t)
}

func TestLessonHandler_GetModuleLessons_Success(t *testing.T) {
	mockLessonRepo := new(MockLessonRepository)
	mockModuleRepo := new(MockModuleRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	lessonUC := usecase.NewLessonUseCase(mockLessonRepo, mockModuleRepo, mockCourseRepo, mockAuditRepo)
	handler := NewLessonHandler(lessonUC)

	moduleID := uuid.New()
	module := &entity.Module{ID: moduleID, Title: "Test Module"}
	lessons := []*entity.LessonVersion{
		{ID: uuid.New(), ModuleID: moduleID, Content: "Lesson 1", VersionNumber: 1},
		{ID: uuid.New(), ModuleID: moduleID, Content: "Lesson 2", VersionNumber: 1},
	}

	mockModuleRepo.On("GetByID", mock.Anything, moduleID).Return(module, nil)
	mockLessonRepo.On("GetLatestByModule", mock.Anything, moduleID).Return(lessons, nil)

	req := httptest.NewRequest(http.MethodGet, "/modules/"+moduleID.String()+"/lessons", nil)
	req = mux.SetURLVars(req, map[string]string{"moduleId": moduleID.String()})
	rr := httptest.NewRecorder()

	handler.GetModuleLessons(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockLessonRepo.AssertExpectations(t)
}

func TestLessonHandler_GetAllLessonVersions_Success(t *testing.T) {
	mockLessonRepo := new(MockLessonRepository)
	mockModuleRepo := new(MockModuleRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	lessonUC := usecase.NewLessonUseCase(mockLessonRepo, mockModuleRepo, mockCourseRepo, mockAuditRepo)
	handler := NewLessonHandler(lessonUC)

	lessonID := uuid.New()
	versions := []*entity.LessonVersion{
		{ID: lessonID, ModuleID: uuid.New(), Content: "Version 1", VersionNumber: 1},
		{ID: lessonID, ModuleID: uuid.New(), Content: "Version 2", VersionNumber: 2},
	}

	mockLessonRepo.On("GetAllVersions", mock.Anything, lessonID).Return(versions, nil)

	req := httptest.NewRequest(http.MethodGet, "/lessons/"+lessonID.String()+"/all-versions", nil)
	req = mux.SetURLVars(req, map[string]string{"lessonId": lessonID.String()})
	rr := httptest.NewRecorder()

	handler.GetAllLessonVersions(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockLessonRepo.AssertExpectations(t)
}

func TestLessonHandler_DeleteLesson_Success(t *testing.T) {
	mockLessonRepo := new(MockLessonRepository)
	mockModuleRepo := new(MockModuleRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	lessonUC := usecase.NewLessonUseCase(mockLessonRepo, mockModuleRepo, mockCourseRepo, mockAuditRepo)
	handler := NewLessonHandler(lessonUC)

	lessonID := uuid.New()
	instructorID := uuid.New()
	moduleID := uuid.New()
	courseID := uuid.New()
	lesson := &entity.LessonVersion{ID: lessonID, ModuleID: moduleID, Content: "Content"}
	module := &entity.Module{ID: moduleID, CourseID: courseID}
	course := &entity.Course{ID: courseID, InstructorID: instructorID}

	mockLessonRepo.On("GetByID", mock.Anything, lessonID).Return(lesson, nil)
	mockModuleRepo.On("GetByID", mock.Anything, moduleID).Return(module, nil)
	mockCourseRepo.On("GetByID", mock.Anything, courseID).Return(course, nil)
	mockLessonRepo.On("Delete", mock.Anything, lessonID).Return(nil)
	mockAuditRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.AuditLog")).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/lessons/"+lessonID.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"lessonId": lessonID.String()})

	claims := &auth.Claims{
		UserID: instructorID,
		Email:  "instructor@example.com",
		Role:   entity.UserRoleInstructor,
	}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, claims)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.DeleteLesson(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
	mockLessonRepo.AssertExpectations(t)
}

func TestLessonHandler_CreateLessonVersion_Success(t *testing.T) {
	mockLessonRepo := new(MockLessonRepository)
	mockModuleRepo := new(MockModuleRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	lessonUC := usecase.NewLessonUseCase(mockLessonRepo, mockModuleRepo, mockCourseRepo, mockAuditRepo)
	handler := NewLessonHandler(lessonUC)

	lessonID := uuid.New()
	moduleID := uuid.New()
	courseID := uuid.New()
	instructorID := uuid.New()

	existingLesson := &entity.LessonVersion{
		ID:            lessonID,
		ModuleID:      moduleID,
		Title:         "Original Lesson",
		VersionNumber: 1,
		Content:       "Original content",
	}
	module := &entity.Module{ID: moduleID, CourseID: courseID}
	course := &entity.Course{ID: courseID, InstructorID: instructorID}

	mockLessonRepo.On("GetByID", mock.Anything, lessonID).Return(existingLesson, nil)
	mockModuleRepo.On("GetByID", mock.Anything, moduleID).Return(module, nil)
	mockCourseRepo.On("GetByID", mock.Anything, courseID).Return(course, nil)
	mockLessonRepo.On("GetNextVersionNumber", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(2, nil)
	mockLessonRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.LessonVersion")).Return(nil)
	mockAuditRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.AuditLog")).Return(nil)

	reqBody := dto.CreateLessonVersionRequest{
		Content:  "Updated content",
		VideoURL: "https://example.com/new-video.mp4",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/lessons/"+lessonID.String()+"/versions", bytes.NewBuffer(body))
	req = mux.SetURLVars(req, map[string]string{"lessonId": lessonID.String()})

	claims := &auth.Claims{
		UserID: instructorID,
		Email:  "instructor@example.com",
		Role:   entity.UserRoleInstructor,
	}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, claims)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.CreateLessonVersion(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	mockLessonRepo.AssertExpectations(t)
	mockModuleRepo.AssertExpectations(t)
	mockCourseRepo.AssertExpectations(t)
	mockAuditRepo.AssertExpectations(t)
}

func TestLessonHandler_CreateLessonVersion_LessonNotFound(t *testing.T) {
	mockLessonRepo := new(MockLessonRepository)
	mockModuleRepo := new(MockModuleRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	lessonUC := usecase.NewLessonUseCase(mockLessonRepo, mockModuleRepo, mockCourseRepo, mockAuditRepo)
	handler := NewLessonHandler(lessonUC)

	lessonID := uuid.New()
	instructorID := uuid.New()

	mockLessonRepo.On("GetByID", mock.Anything, lessonID).Return(nil, entity.ErrNotFound)

	reqBody := dto.CreateLessonVersionRequest{
		Content:  "Updated content",
		VideoURL: "https://example.com/new-video.mp4",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/lessons/"+lessonID.String()+"/versions", bytes.NewBuffer(body))
	req = mux.SetURLVars(req, map[string]string{"lessonId": lessonID.String()})

	claims := &auth.Claims{
		UserID: instructorID,
		Email:  "instructor@example.com",
		Role:   entity.UserRoleInstructor,
	}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, claims)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.CreateLessonVersion(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	mockLessonRepo.AssertExpectations(t)
}

func TestLessonHandler_CreateLessonVersion_Unauthorized(t *testing.T) {
	mockLessonRepo := new(MockLessonRepository)
	mockModuleRepo := new(MockModuleRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	lessonUC := usecase.NewLessonUseCase(mockLessonRepo, mockModuleRepo, mockCourseRepo, mockAuditRepo)
	handler := NewLessonHandler(lessonUC)

	lessonID := uuid.New()
	moduleID := uuid.New()
	courseID := uuid.New()
	instructorID := uuid.New()
	otherUserID := uuid.New()

	existingLesson := &entity.LessonVersion{
		ID:            lessonID,
		ModuleID:      moduleID,
		Title:         "Original Lesson",
		VersionNumber: 1,
		Content:       "Original content",
	}
	module := &entity.Module{ID: moduleID, CourseID: courseID}
	course := &entity.Course{ID: courseID, InstructorID: instructorID}

	mockLessonRepo.On("GetByID", mock.Anything, lessonID).Return(existingLesson, nil)
	mockModuleRepo.On("GetByID", mock.Anything, moduleID).Return(module, nil)
	mockCourseRepo.On("GetByID", mock.Anything, courseID).Return(course, nil)

	reqBody := dto.CreateLessonVersionRequest{
		Content:  "Updated content",
		VideoURL: "https://example.com/new-video.mp4",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/lessons/"+lessonID.String()+"/versions", bytes.NewBuffer(body))
	req = mux.SetURLVars(req, map[string]string{"lessonId": lessonID.String()})

	claims := &auth.Claims{
		UserID: otherUserID,
		Email:  "other@example.com",
		Role:   entity.UserRoleInstructor,
	}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, claims)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.CreateLessonVersion(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
	mockLessonRepo.AssertExpectations(t)
	mockModuleRepo.AssertExpectations(t)
	mockCourseRepo.AssertExpectations(t)
}
