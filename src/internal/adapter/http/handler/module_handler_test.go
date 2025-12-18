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

func TestModuleHandler_CreateModule_Success(t *testing.T) {
	mockModuleRepo := new(MockModuleRepository)
	mockCourseRepo := new(MockCourseRepository)

	moduleUC := usecase.NewModuleUseCase(mockModuleRepo, mockCourseRepo)
	handler := NewModuleHandler(moduleUC)

	courseID := uuid.New().String()
	instructorID := uuid.New().String()
	course := &entity.Course{ID: courseID, InstructorID: instructorID, Title: "Test Course"}

	mockCourseRepo.On("GetByID", mock.Anything, courseID).Return(course, nil)
	mockModuleRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Module")).Return(nil)

	reqBody := dto.CreateModuleRequest{Title: "Module 1", OrderIndex: 0}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/courses/"+courseID+"/modules", bytes.NewBuffer(body))
	req = mux.SetURLVars(req, map[string]string{"courseId": courseID})

	claims := &auth.Claims{UserID: instructorID, Email: "instructor@example.com", Role: entity.UserRoleInstructor}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, claims)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.CreateModule(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	mockModuleRepo.AssertExpectations(t)
	mockCourseRepo.AssertExpectations(t)
}

func TestModuleHandler_GetCourseModules_Success(t *testing.T) {
	mockModuleRepo := new(MockModuleRepository)
	mockCourseRepo := new(MockCourseRepository)

	moduleUC := usecase.NewModuleUseCase(mockModuleRepo, mockCourseRepo)
	handler := NewModuleHandler(moduleUC)

	courseID := uuid.New().String()
	course := &entity.Course{ID: courseID, Title: "Test Course"}
	modules := []*entity.Module{
		{ID: uuid.New().String(), CourseID: courseID, Title: "Module 1"},
		{ID: uuid.New().String(), CourseID: courseID, Title: "Module 2"},
	}

	mockCourseRepo.On("GetByID", mock.Anything, courseID).Return(course, nil)
	mockModuleRepo.On("GetByCourse", mock.Anything, courseID).Return(modules, nil)

	req := httptest.NewRequest(http.MethodGet, "/courses/"+courseID+"/modules", nil)
	req = mux.SetURLVars(req, map[string]string{"courseId": courseID})
	rr := httptest.NewRecorder()

	handler.GetCourseModules(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockModuleRepo.AssertExpectations(t)
}

func TestModuleHandler_GetModule_Success(t *testing.T) {
	mockModuleRepo := new(MockModuleRepository)
	mockCourseRepo := new(MockCourseRepository)

	moduleUC := usecase.NewModuleUseCase(mockModuleRepo, mockCourseRepo)
	handler := NewModuleHandler(moduleUC)

	moduleID := uuid.New().String()
	module := &entity.Module{ID: moduleID, Title: "Module 1", CourseID: uuid.New().String()}

	mockModuleRepo.On("GetByID", mock.Anything, moduleID).Return(module, nil)

	req := httptest.NewRequest(http.MethodGet, "/modules/"+moduleID, nil)
	req = mux.SetURLVars(req, map[string]string{"id": moduleID})
	rr := httptest.NewRecorder()

	handler.GetModule(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockModuleRepo.AssertExpectations(t)
}

func TestModuleHandler_DeleteModule_Success(t *testing.T) {
	mockModuleRepo := new(MockModuleRepository)
	mockCourseRepo := new(MockCourseRepository)

	moduleUC := usecase.NewModuleUseCase(mockModuleRepo, mockCourseRepo)
	handler := NewModuleHandler(moduleUC)

	moduleID := uuid.New().String()
	instructorID := uuid.New().String()
	courseID := uuid.New().String()
	module := &entity.Module{ID: moduleID, CourseID: courseID, Title: "Module 1"}
	course := &entity.Course{ID: courseID, InstructorID: instructorID}

	mockModuleRepo.On("GetByID", mock.Anything, moduleID).Return(module, nil)
	mockCourseRepo.On("GetByID", mock.Anything, courseID).Return(course, nil)
	mockModuleRepo.On("Delete", mock.Anything, moduleID).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/modules/"+moduleID, nil)
	req = mux.SetURLVars(req, map[string]string{"id": moduleID})

	claims := &auth.Claims{UserID: instructorID, Email: "instructor@example.com", Role: entity.UserRoleInstructor}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, claims)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.DeleteModule(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
	mockModuleRepo.AssertExpectations(t)
}
